package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/mechta-market/gotemplate/internal/config"
	"github.com/mechta-market/gotemplate/pkg/proto/common"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type GrpcGwServer struct {
	server *http.Server
}

type HttpRoute struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request, map[string]string)
}

func NewGrpcGwServer(
	handlers []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error,
	httpRoutes []HttpRoute,
) (*GrpcGwServer, error) {
	handler, err := GrpcGatewayCreateHandler(func(mux *runtime.ServeMux) error {
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		conn, err := grpc.NewClient("localhost:"+config.Conf.GrpcPort, opts...)
		if err != nil {
			return fmt.Errorf("grpc.NewClient: %w", err)
		}

		for _, h := range handlers {
			err = h(context.Background(), mux, conn)
			if err != nil {
				return fmt.Errorf("grpc-gateway: register grpc-handler: %w", err)
			}
		}

		for _, r := range httpRoutes {
			if err := mux.HandlePath(r.Method, r.Path, r.Handler); err != nil {
				return fmt.Errorf("mux.HandlePath: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("GrpcGatewayCreateHandler: %w", err)
	}

	// server
	server := &http.Server{
		Addr:              ":" + config.Conf.HttpPort,
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       time.Minute,
		MaxHeaderBytes:    300 * 1024,
	}

	return &GrpcGwServer{
		server: server,
	}, nil
}

func (s *GrpcGwServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http-server stopped", "error", err)
			os.Exit(1)
		}
	}()
	slog.Info("http-server started " + s.server.Addr)
}

func (s *GrpcGwServer) Stop(ctx context.Context) error {
	ctx, ctxCancel := context.WithTimeout(ctx, 15*time.Second)
	defer ctxCancel()

	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("http-server shutdown error", "error", err)
		return err
	}

	return nil
}

func GrpcGatewayCreateHandler(muxHook func(*runtime.ServeMux) error) (http.Handler, error) {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithErrorHandler(func(_ context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			var repBody []byte
			var repObj common.ErrorRep
			var headerStatus int

			if st, ok := status.FromError(err); ok {
				headerStatus = runtime.HTTPStatusFromCode(st.Code())
				repObj.Fields = map[string]string{}
				if len(st.Details()) > 0 {
					if detail, ok := st.Details()[0].(*common.ErrorRep); ok {
						repObj.Code = detail.Code
						repObj.Message = detail.Message
						maps.Copy(repObj.Fields, detail.Fields)
					}
				} else {
					switch st.Code() {
					case codes.NotFound:
						repObj.Message = "service path not found"
					case codes.Unimplemented:
						headerStatus = http.StatusMethodNotAllowed
						repObj.Message = "method not allowed"
					case codes.PermissionDenied:
						repObj.Message = "permission denied"
					case codes.Unavailable, codes.Internal:
						repObj.Message = "service not available"
					default:
						headerStatus = http.StatusBadRequest
						repObj.Message = err.Error()
						//repObj.Message = fmt.Sprintf("%s: %v", "invalid argument", st.Message())
					}
				}
			} else {
				headerStatus = http.StatusBadRequest
				repObj.Message = err.Error()
			}

			if len(repBody) == 0 {
				repBody, err = protojson.Marshal(&repObj)
				if err != nil {
					slog.Error("GRPC_GW: ErrorHandler: Failed to marshal", "error", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(headerStatus)
			_, err = io.Copy(w, bytes.NewReader(repBody))
			if err != nil {
				slog.Error("GRPC_GW: ErrorHandler: Failed to write response", "error", err)
			}
		}))

	if muxHook != nil {
		err := muxHook(mux)
		if err != nil {
			return nil, fmt.Errorf("grpc-gateway: muxHook: %w", err)
		}
	}

	// add health check handler
	err := mux.HandlePath(http.MethodGet, "/healthcheck", func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
	})
	if err != nil {
		return nil, fmt.Errorf("grpc-gateway: register healthcheck handler: %w", err)
	}

	// add docs handler
	docFS := http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs")))
	err = mux.HandlePath(http.MethodGet, "/docs/{any}", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		docFS.ServeHTTP(w, r)
	})
	if err != nil {
		return nil, fmt.Errorf("grpc-gateway: register docs handler: %w", err)
	}
	err = mux.HandlePath(http.MethodGet, "/docs/{any}/{any}", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		docFS.ServeHTTP(w, r)
	})
	if err != nil {
		return nil, fmt.Errorf("grpc-gateway: register docs handler: %w", err)
	}

	// add metrics handler
	if err = mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}); err != nil {
		return nil, fmt.Errorf("grpc-gateway: register metrics handler: %w", err)
	}

	handler := http.Handler(mux)

	// add cors middleware
	if config.Conf.HttpCors {
		handler = cors.New(cors.Options{
			AllowOriginFunc: func(_ string) bool { return true },
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPut,
				http.MethodPost,
				http.MethodDelete,
			},
			AllowedHeaders: []string{
				"Accept",
				"Content-Type",
				"X-Requested-With",
				"Authorization",
			},
			AllowCredentials: true,
			MaxAge:           604800,
		}).Handler(handler)
	}

	// add recover middleware
	handler = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				// use always new err instance in defer
				if err := recover(); err != nil {
					slog.Error(
						"Recovered from panic",
						slog.Any("error", err),
						slog.Any("recovery_stacktrace", string(debug.Stack())),
					)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}(handler)

	return handler, nil
}
