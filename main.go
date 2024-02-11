package main

import "math"

type NeuralNetwork struct {
	neurons [][]float64   // оутпуты на нейроне и сами нейроны
	w       [][][]float64 // веса [слой][нейрон][вес]
	LR      float64       // скорость обучения
	EPOCH   int           // кол-во эпох обучения
}

func main() {

}

func loadDataSet() (trainData [][]float64, expResults []float64) { // загрузка data set'а из файла
	return nil, nil
}

func loadTestData() [][]float64 { // загрузка тестового data set'а из файла
	return nil
}

func saveWeights() {} // сохранить мозги (веса) в файл

func loadWeights() {} // загрузить пресет мозгов (весов)

func activate(s float64) float64 { // функция активации
	return math.Tanh(s)
}

func forward(inputs []float64, w [][]float64) float64 { // функция прямого распространения
	return 0
}

func backProp() { // Функция обратного распространения ошибки

}

func train(nn *NeuralNetwork, data [][]float64, exp []float64) {

	for e := 0; e < nn.EPOCH; e++ { // цикл по эпохам

		for d := 0; d < len(data); d++ { // цикл по дата сету // вынести в функцию

			for n := 0; n < len(nn.w); n++ { // цикл по слоям нейронов
				// прогоняем все слои нейронов
				// В цикле прогнать forward для каждого нейрона в слое, можно цикл сделать в forward()

			}

			// вычисляем ошибку для каждого слоя
			// меняем веса

		}

	}
}
