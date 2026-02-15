package trainings

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Training хранит все данные о тренировке.
// Steps — количество шагов за тренировку.
// TrainingType — тип тренировки ("Бег" или "Ходьба").
// Duration — длительность тренировки.
// personaldata.Personal — встроенная (эмбеддинг) структура с данными пользователя.
// Благодаря эмбеддингу мы получаем прямой доступ к полям Name, Weight, Height
// и к методу Print() — без необходимости писать t.Personal.Weight, можно просто t.Weight.
type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

// Parse разбирает строку формата "3456,Ходьба,3h00m" и записывает данные
// в соответствующие поля структуры Training.
// Используем указатель (*Training), потому что мы записываем данные в структуру —
// без указателя изменения бы не сохранились (Go передаёт аргументы по значению).
//
// ИСПРАВЛЕНО:
// Было (оригинальный код):
//
//	func (t *Training) Parse(datastring string) error {
//	    parts := strings.Split(datastring, ",")
//	    if len(parts) != 3 {
//	        return fmt.Errorf("invalid data format")
//	    }
//	    steps, err := strconv.Atoi(parts[0])
//	    if err != nil {
//	        return err
//	    }
//	    duration, err := time.ParseDuration(parts[2])
//	    if err != nil {
//	        return err
//	    }
//	    t.Steps = steps            // <-- записывались без проверки!
//	    t.TrainingType = parts[1]  // <-- записывались без проверки!
//	    t.Duration = duration      // <-- записывались без проверки!
//	    return nil
//	}
//
// Проблема: не было проверки steps <= 0 и duration <= 0.
// Тесты ожидают ошибку при нулевых и отрицательных шагах/длительности,
// а старый код просто записывал их в структуру без валидации.
//
// Стало: добавлены две проверки (if steps <= 0 и if duration <= 0)
// перед записью в структуру.
//
// Также была УДАЛЕНА неиспользуемая функция NewTraining(), которая
// была в оригинале, но не нужна — структуру можно создать напрямую
// через Training{Personal: person} (так и делается в main.go).
func (t *Training) Parse(datastring string) error {
	// Разделяем строку по запятой. Ожидаем ровно 3 части:
	// [0] — количество шагов, [1] — тип тренировки, [2] — длительность.
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		// Если частей не 3, формат строки некорректный — возвращаем ошибку.
		return fmt.Errorf("invalid data format")
	}

	// Преобразуем первый элемент слайса (строку) в int.
	// strconv.Atoi вернёт ошибку, если строка не является числом
	// (например "abc" или пустая строка).
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}

	// ДОБАВЛЕНО: проверяем, что количество шагов положительное.
	// Ноль или отрицательное значение не имеют смысла для тренировки.
	if steps <= 0 {
		return fmt.Errorf("шаги должны быть положительным числом")
	}

	// Парсим строку длительности (например "3h00m") в time.Duration.
	// time.ParseDuration — стандартная функция Go для парсинга строк
	// в формате "1h30m", "45m", "2h" и т.д.
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return err
	}

	// ДОБАВЛЕНО: проверяем, что длительность положительная.
	// Нулевая или отрицательная длительность не имеет смысла.
	if duration <= 0 {
		return fmt.Errorf("длительность должна быть положительной")
	}

	// Все проверки пройдены — сохраняем данные в полях структуры.
	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration
	return nil
}

// ActionInfo формирует и возвращает строку с данными о тренировке.
// Метод не использует указатель (t Training), потому что только читает данные.
// Возвращает строку с результатами и ошибку (например, если тип тренировки неизвестен).
// Этот метод был реализован правильно в оригинале — изменений не потребовалось.
func (t Training) ActionInfo() (string, error) {
	// Вычисляем дистанцию: шаги * (рост * 0.45) / 1000.
	distance := spentenergy.Distance(t.Steps, t.Height)
	// Вычисляем среднюю скорость: дистанция / время в часах.
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var calories float64
	var err error

	// В зависимости от типа тренировки используем разные формулы расчёта калорий.
	// switch — удобнее цепочки if-else, когда сравниваем одну переменную с несколькими значениями.
	switch t.TrainingType {
	case "Бег":
		// Формула бега: (вес * скорость * минуты) / 60
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Ходьба":
		// Формула ходьбы: ((вес * скорость * минуты) / 60) * 0.5
		// Коэффициент 0.5, потому что при ходьбе тратится меньше энергии.
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		// Если передан неизвестный тип тренировки — возвращаем ошибку.
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	// Формируем итоговую строку с информацией.
	// Duration.Hours() переводит длительность в часы (дробное число).
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
