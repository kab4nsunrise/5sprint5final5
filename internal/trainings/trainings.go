package trainings

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}


func (t *Training) Parse(datastring string) error {
	
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		
		return fmt.Errorf("invalid data format")
	}

	
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}

	
	if steps <= 0 {
		return fmt.Errorf("шаги должны быть положительным числом")
	}

	
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return err
	}

	
	if duration <= 0 {
		return fmt.Errorf("длительность должна быть положительной")
	}

	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration
	return nil
}

func (t Training) ActionInfo() (string, error) {

	distance := spentenergy.Distance(t.Steps, t.Height)
	
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var calories float64
	var err error

	
	switch t.TrainingType {
	case "Бег":
		
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Ходьба":
		
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		t.TrainingType,
		t.Duration.Hours(),
		distance,
		speed,
		calories,
	), nil
}
