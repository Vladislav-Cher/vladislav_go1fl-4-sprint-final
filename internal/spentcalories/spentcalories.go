package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	dataInSlice := strings.Split(data, ",")
	if len(dataInSlice) != 3 {
		return 0, "", time.Duration(0), fmt.Errorf("требуется 3 элемента данных (получено %d)", len(dataInSlice))
	}

	//Cделал раздельно по две проверки для каждой переменной,
	//чтобы в случае нулевых значений не возвращалсь нули с "nil" ошибкой.
	steps, err := strconv.Atoi(dataInSlice[0])
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	if steps <= 0 {
		return 0, "", time.Duration(0), fmt.Errorf("некорректное количество шагов (получено %d)", steps)
	}

	duration, err := time.ParseDuration(dataInSlice[2])
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	if duration <= 0 {
		return 0, "", time.Duration(0), fmt.Errorf("некорректная продолжительность тренировки")
	}

	return steps, dataInSlice[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	distInM := float64(steps) * stepLength
	distInKm := distInM / float64(mInKm)

	return distInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if weight <= 0 {
		return "", fmt.Errorf("некорректный вес (получено %.1f кг)", weight)
	}

	if height <= 0 {
		return "", fmt.Errorf("некорректный рост (получено %.2f м)", height)
	}

	var dist float64
	var meanS float64
	var calories float64
	err = nil

	switch activity {
	case "Бег":
		dist = distance(steps, height)
		meanS = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		dist = distance(steps, height)
		meanS = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км."+
			"\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		dist,
		meanS,
		calories,
	), nil

}

func variablesCheck(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 {
		return fmt.Errorf("некорректное количество шагов (получено %d)", steps)
	}
	if weight <= 0 {
		return fmt.Errorf("некорректный вес (получено %.1f кг)", weight)
	}
	if height <= 0 {
		return fmt.Errorf("некорректный рост (получено %.2f м)", height)
	}
	if duration <= 0 {
		return fmt.Errorf("некорректная продолжительность тренировки")
	}
	return nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	variablesCheck(steps, weight, height, duration)

	meanS := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	return (weight * meanS * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	variablesCheck(steps, weight, height, duration)

	meanS := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	caloriesIfRunning := (weight * meanS * durationInMinutes) / minInH

	return caloriesIfRunning * walkingCaloriesCoefficient, nil
}
