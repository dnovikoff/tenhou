package utils

type Task interface {
	// Called on one of worker routine
	Run() error
	// Called on single routine
	Save(blocked bool) error
}

type Scheduler struct {
	routines   Routines
	bgRoutines Routines
	tasks      chan Task
	save       chan Task
}

func (s *Scheduler) Start(cnt, saveSize int) {
	s.tasks = make(chan Task, cnt)
	s.save = make(chan Task, saveSize)
	for i := 0; i < cnt; i++ {
		s.routines.Start(func() error {
			for t := range s.tasks {
				err := t.Run()
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	s.bgRoutines.Start(s.saveRoutine)
}

func (s *Scheduler) saveRoutine() error {
	for {
		t, blocked := s.readForSave()
		if t == nil {
			return nil
		}
		err := t.Save(blocked)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scheduler) readForSave() (Task, bool) {
	var t Task
	blocked := false
	select {
	case t = <-s.save:
	default:
		t = <-s.save
		blocked = true
	}
	return t, blocked
}

func (s *Scheduler) Stop() error {
	close(s.tasks)
	err := s.routines.Wait()
	close(s.save)
	bgErr := s.bgRoutines.Wait()
	if err != nil {
		return err
	}
	return bgErr
}

func (s *Scheduler) Push(t Task) {
	s.tasks <- t
}

func (s *Scheduler) Error() error {
	err := s.routines.Error()
	if err != nil {
		return err
	}
	return s.bgRoutines.Error()
}
