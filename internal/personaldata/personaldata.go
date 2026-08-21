package personaldata

import "fmt"

type Personal struct {
	Name   string  // Имя пользователя
	Weight float64 // Вес пользователя в килограммах
	Height float64 // Рост пользователя в метрах
}


func (p Personal) Print() {
	fmt.Printf(
		"Имя: %s\nВес: %.2f кг.\nРост: %.2f м.\n",
		p.Name,
		p.Weight,
		p.Height,
	)
}
