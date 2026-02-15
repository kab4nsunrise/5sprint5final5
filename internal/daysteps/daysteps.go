package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}


func (ds *DaySteps) Parse(datastring string) (err error) {
	
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		
		return fmt.Errorf("invalid data format")
	}

	
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}

	
	if steps <= 0 {
		return fmt.Errorf("шаги должны быть положительным числом")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return err
	}
	if duration <= 0 {
		return fmt.Errorf("длительность должна быть положительной")
	}


	return nil
}


	calories, err := spentenergy.WalkingSpentCalories(
		ds.Steps,
		ds.Weight,
		ds.Height,
		ds.Duration,
	)
	if err != nil {
	
		return "", err
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		ds.Steps,
		distance,
		calories,
	), nil
}
