package main

import (
	"fmt"
	"math"
)

type NeuralNetwork struct {
	neurons [][]float64   // оутпуты на нейроне и сами нейроны
	w       [][][]float64 // веса [слой][нейрон][вес]
	LR      float64       // скорость обучения
	EPOCH   int           // кол-во эпох обучения
}

func NewNeuralNetwork(neurons [][]float64, w [][][]float64, LR float64, EPOCH int) *NeuralNetwork {
	return &NeuralNetwork{neurons: neurons, w: w, LR: LR, EPOCH: EPOCH}
}

func createNeurons(a ...int) (neurons [][]float64) {
	neurons = make([][]float64, len(a))
	for i, c := range a {
		neurons[i] = make([]float64, c)
	}
	return neurons
}

func createWeights(neurons [][]float64) (w [][][]float64) {
	w = make([][][]float64, len(neurons)-1)
	for i := 0; i < len(neurons)-1; i++ {
		w[i] = make([][]float64, len(neurons[i]))
		for j := 0; j < len(neurons[i]); j++ {
			w[i][j] = make([]float64, len(neurons[i+1]))
		}
	}
	return w
}

const a = 1 // a - альфа для функции активации

func main() {
	nn := createNeurons(2, 3, 2)
	w := createWeights(nn)
	fmt.Println("%v", nn)
	fmt.Println("%v", w)
	//var nn = NeuralNetwork{neurons: make([][]float64, 5), w: make([][][]float64, 4)}
	//nn.neurons[0] = []float64{0.5, 0.2, 0, 1}
	//nn.neurons[1] = []float64{0.3, 0.1, 1, 1}
	//nn.neurons[2] = []float64{0.3, 0.1, 1, 1}
	//nn.neurons[3] = []float64{0.3, 0.1, 1, 1}
	//nn.neurons[4] = []float64{0.3, 0.1, 1, 1}
	//nn.w[0] = [][]float64{{1, 2, 1, 0.5}, {1, 2, 0, 0.5}, {1, 1, 1, 0.5}, {1, 0.2, 1, 0.2}}
	//nn.w[1] = [][]float64{{1, 2, 1, 0.5}, {1, 2, 0, 0.5}, {1, 1, 1, 0.5}, {1, 0.2, 1, 0.2}}
	//nn.w[2] = [][]float64{{1, 2, 1, 0.5}, {1, 2, 0, 0.5}, {1, 1, 1, 0.5}, {1, 0.2, 1, 0.2}}
	//nn.w[3] = [][]float64{{1, 2, 1, 0.5}, {1, 2, 0, 0.5}, {1, 1, 1, 0.5}, {1, 0.2, 1, 0.2}}
	//
	//fmt.Printf("NN:\n %v\n", nn.neurons)
	//fmt.Printf("W:\n %v\n", nn.w)
	//
	////forward(nn.neurons[0], nn.neurons[1], nn.w[0])
	//
	//for n := 0; n < len(nn.w); n++ { // цикл по слоям нейронов
	//	// прогоняем все слои нейронов
	//	// В цикле прогнать forward для каждого нейрона в слое, можно цикл сделать в forward()
	//	forward(nn.neurons[n], nn.neurons[n+1], nn.w[n]) // текущий слой, следующий слой, веса между этими слоями
	//	fmt.Printf("iteration n: %d \n NN: %v\n", n, nn.neurons)
	//}
	//
	//fmt.Printf("NN:\n %v\n", nn.neurons)
	//fmt.Printf("W:\n %v\n", nn.w)
}

func loadDataSet() (trainData [][]float64, expResults []float64) { // загрузка data set'а из файла
	return nil, nil
}

func loadTestData() [][]float64 { // загрузка тестового data set'а из файла
	return nil
}

func saveWeights() {} // сохранить мозги (веса) в файл

func loadWeights() {} // загрузить пресет мозгов (весов)

func generateWeights(w [][][]float64) {} // сгенерировать рандомно веса

func activate(s float64) float64 { // функция активации
	return math.Tanh(a * s) // гиперболический тангенс
	//return 1 / (1 + math.Pow(math.E, -(a*s))) // сигмоидная функция
}

func forward(inputNeurons []float64, outputNeurons []float64, w [][]float64) { // функция прямого распространения
	for i, _ := range outputNeurons {
		var res float64 = 0
		for j, input := range inputNeurons {
			res += input * w[j][i] // i - к output нейрону, j - из input нейрона
		}
		outputNeurons[i] = activate(res)
	}
}

func backProp() { // Функция обратного распространения ошибки
	// должен тут вызывать forward() и что-то получать и дальше считать цену ошибки

}

func train(nn *NeuralNetwork, data [][]float64, exp []float64) {

	for e := 0; e < nn.EPOCH; e++ { // цикл по эпохам

		for d := 0; d < len(data); d++ { // цикл по дата сету // вынести в функцию

			for n := 0; n < len(nn.w); n++ { // цикл по слоям нейронов
				// прогоняем все слои нейронов
				// В цикле прогнать forward для каждого нейрона в слое, можно цикл сделать в forward()
				forward(nn.neurons[n], nn.neurons[n+1], nn.w[n]) // текущий слой, следующий слой, веса между этими слоями
				//fmt.Printf("%v", nn.neurons)
			}

			// вычисляем ошибку для каждого слоя
			// меняем веса

		}

	}
}
