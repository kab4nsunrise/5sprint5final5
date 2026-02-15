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
		// Если частей не 2, формат строки некорректный.
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

	// Все проверки пройдены — записываем данные в поля структуры.
	ds.Steps = steps
	ds.Duration = duration
	return nil
}

func (ds DaySteps) ActionInfo() (string, error) {
	
	distance := spentenergy.Distance(ds.Steps, ds.Height)

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
