package example

type Usecase struct {
	service iService
}

func New(
	service iService,
) *Usecase {
	return &Usecase{
		service: service,
	}
}

//func (u *Usecase) List
//func (u *Usecase) Get
//func (u *Usecase) Set
//func (u *Usecase) Delete
