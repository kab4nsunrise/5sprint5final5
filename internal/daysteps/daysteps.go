package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// DaySteps хранит данные о дневной прогулке.
// Steps — количество шагов, Duration — длительность прогулки.
// personaldata.Personal — встроенная структура (эмбеддинг), благодаря
// которой у нас есть доступ к полям Name, Weight, Height и методу Print()
// без необходимости обращаться через отдельное поле.
type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

// Parse разбирает строку формата "678,0h50m" и записывает данные
// в поля структуры DaySteps.
// Используем указатель на структуру (*DaySteps), потому что мы изменяем
// её поля — без указателя изменения не сохранились бы (Go передаёт по значению).
//
// ИСПРАВЛЕНО:
// Было (оригинальный код):
//
//	func (ds *DaySteps) Parse(datastring string) error {
//	    parts := strings.Split(datastring, ",")
//	    if len(parts) != 2 {
//	        return fmt.Errorf("invalid data format")
//	    }
//	    steps, err := strconv.Atoi(parts[0])
//	    if err != nil {
//	        return err
//	    }
//	    duration, err := time.ParseDuration(parts[1])
//	    if err != nil {
//	        return err
//	    }
//	    ds.Steps = steps        // <-- записывались без проверки!
//	    ds.Duration = duration   // <-- записывались без проверки!
//	    return nil
//	}
//
// Проблема: не было проверки steps <= 0 и duration <= 0.
// Тесты ожидают ошибку при нулевых и отрицательных шагах/длительности,
// а старый код просто записывал их в структуру без валидации.
//
// Стало: добавлены две проверки (if steps <= 0 и if duration <= 0)
// перед записью в структуру. Только после прохождения всех проверок
// данные сохраняются в полях.
func (ds *DaySteps) Parse(datastring string) (err error) {
	// Разделяем строку по запятой. Ожидаем ровно 2 части:
	// [0] — количество шагов, [1] — длительность.
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		// Если частей не 2, формат строки некорректный.
		return fmt.Errorf("invalid data format")
	}

	// Преобразуем первую часть (строку) в целое число (количество шагов).
	// strconv.Atoi вернёт ошибку, если строка не является числом.
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}

	// ДОБАВЛЕНО: проверяем, что количество шагов положительное.
	// Ноль или отрицательное значение не имеют смысла для прогулки.
	if steps <= 0 {
		return fmt.Errorf("шаги должны быть положительным числом")
	}

	// Парсим строку длительности (например "0h50m") в time.Duration.
	// time.ParseDuration понимает форматы вроде "1h30m", "45m", "2h" и т.д.
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return err
	}

	// ДОБАВЛЕНО: проверяем, что длительность положительная.
	// Нулевая или отрицательная длительность не имеет смысла.
	if duration <= 0 {
		return fmt.Errorf("длительность должна быть положительной")
	}

	// Все проверки пройдены — записываем данные в поля структуры.
	ds.Steps = steps
	ds.Duration = duration
	return nil
}

// ActionInfo формирует и возвращает строку с информацией о прогулке.
// Метод не использует указатель (ds DaySteps), потому что он только читает
// данные из структуры и ничего в ней не меняет.
// Возвращает строку с результатами и ошибку, если что-то пошло не так.
// Этот метод был реализован правильно в оригинале — изменений не потребовалось.
func (ds DaySteps) ActionInfo() (string, error) {
	// Вычисляем дистанцию в километрах по количеству шагов и росту.
	// Функция Distance из пакета spentenergy: шаги * (рост * 0.45) / 1000.
	distance := spentenergy.Distance(ds.Steps, ds.Height)

	// Вычисляем калории, потраченные при ходьбе.
	// WalkingSpentCalories возвращает ошибку, если входные данные некорректны
	// (например, нулевой вес или рост).
	calories, err := spentenergy.WalkingSpentCalories(
		ds.Steps,
		ds.Weight,
		ds.Height,
		ds.Duration,
	)
	if err != nil {
		// При ошибке возвращаем пустую строку и саму ошибку.
		return "", err
	}

	// Формируем итоговую строку с информацией о прогулке.
	// %.2f — формат с двумя знаками после запятой.
	return fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		ds.Steps,
		distance,
		calories,
	), nil
}
