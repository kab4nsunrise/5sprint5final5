package spentenergy

import (
	"errors"
	"time"
)

const (
	
	stepLengthCoefficient = 0.45
	mInKm = 1000
	minInH = 60

	walkingCaloriesCoefficient = 0.5
)


func Distance(steps int, height float64) float64 {
	
	if steps <= 0 || height <= 0 {
		return 0
	}
		stepLength := height * stepLengthCoefficient

	return float64(steps) * stepLength / mInKm
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	
	if steps <= 0 || duration <= 0 {
		return 0
	}
	
	return Distance(steps, height) / duration.Hours()
}


func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid input data")
	}
	
	speed := MeanSpeed(steps, height, duration)

	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid input data")
	}
	speed := MeanSpeed(steps, height, duration)
	
	calories := (weight * speed * duration.Minutes()) / minInH
	
	return calories * walkingCaloriesCoefficient, nil
}
