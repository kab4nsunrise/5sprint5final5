// Пакет spentenergy содержит функции для расчёта дистанции, средней скорости
// и потраченных калорий при беге и ходьбе.
// Все функции экспортируемые (начинаются с заглавной буквы), потому что они
// используются в других пакетах (daysteps, trainings).
//
// Этот пакет был реализован правильно в оригинале — код не менялся,
// добавлены только поясняющие комментарии.
package spentenergy

import (
	"errors"
	"time"
)

const (
	// stepLengthCoefficient — коэффициент для расчёта длины шага.
	// Длина шага = рост * 0.45. Это приблизительная формула.
	stepLengthCoefficient = 0.45

	// mInKm — количество метров в километре, для перевода единиц.
	mInKm = 1000

	// minInH — количество минут в часе, для нормализации калорий.
	minInH = 60

	// walkingCaloriesCoefficient — коэффициент для ходьбы.
	// При ходьбе тратится примерно вдвое меньше калорий, чем при беге,
	// поэтому результат умножается на 0.5.
	walkingCaloriesCoefficient = 0.5
)

// Distance вычисляет дистанцию в километрах по количеству шагов и росту.
// Формула: (шаги * рост * 0.45) / 1000
// Обратите внимание: steps имеет тип int, поэтому мы приводим его к float64
// через float64(steps), чтобы результат деления был дробным числом.
func Distance(steps int, height float64) float64 {
	// Проверяем входные данные — нулевые и отрицательные значения бессмысленны.
	if steps <= 0 || height <= 0 {
		return 0
	}
	// Длина одного шага = рост * коэффициент.
	stepLength := height * stepLengthCoefficient
	// Дистанция = количество шагов * длина шага, делим на 1000 (переводим метры в км).
	return float64(steps) * stepLength / mInKm
}

// MeanSpeed вычисляет среднюю скорость в км/ч.
// Формула: дистанция / время_в_часах.
// duration.Hours() — метод из пакета time, который переводит Duration в часы (дробное число).
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем: шаги должны быть положительными, длительность тоже.
	// Если duration <= 0, делить на неё нельзя (деление на ноль или бессмыслица).
	if steps <= 0 || duration <= 0 {
		return 0
	}
	// Используем уже написанную функцию Distance для вычисления дистанции.
	return Distance(steps, height) / duration.Hours()
}

// RunningSpentCalories вычисляет калории, потраченные при беге.
// Формула: (вес * средняя_скорость * минуты) / 60
// Возвращает два значения: калории и ошибку.
// Ошибка возвращается, если входные данные некорректны (нулевые или отрицательные).
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем ВСЕ параметры на корректность.
	// Любой из них <= 0 делает расчёт бессмысленным.
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid input data")
	}
	// Вычисляем среднюю скорость с помощью уже написанной функции.
	speed := MeanSpeed(steps, height, duration)
	// duration.Minutes() переводит длительность в минуты (дробное число).
	// Делим на minInH (60), чтобы нормализовать результат к часовой шкале.
	return (weight * speed * duration.Minutes()) / minInH, nil
}

// WalkingSpentCalories вычисляет калории, потраченные при ходьбе.
// Формула такая же, как при беге, но результат умножается на коэффициент 0.5,
// потому что при ходьбе расход энергии ниже, чем при беге.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Та же проверка, что и в RunningSpentCalories.
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid input data")
	}
	speed := MeanSpeed(steps, height, duration)
	// Считаем базовые калории (как при беге).
	calories := (weight * speed * duration.Minutes()) / minInH
	// Умножаем на walkingCaloriesCoefficient (0.5) — корректирующий коэффициент для ходьбы.
	return calories * walkingCaloriesCoefficient, nil
}
