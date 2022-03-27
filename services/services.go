package services

import "gorm.io/gorm"

type Services struct {
	UserService       UserService
	ProgramService    ProgramService
	CourseUnitService CourseUnitService
	ResourceService   ResourceService
}

func New(db *gorm.DB) *Services {
	return &Services{
		UserService:       NewUserService(db),
		ProgramService:    NewProgramService(db),
		CourseUnitService: NewCourseUnitService(db),
		ResourceService:   NewResourceService(db),
	}
}
