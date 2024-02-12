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

//TODO: чтобы учитывалось смещение B biases + просчитывалось в forward() и при создании весов тоже. Можно сделать как Доп.Нейрон

//TODO: сделать обработку ошибок с тем что дата сет и кол-во входных нейронов совпадало
//TODO: сделать обработку ошибок с тем чтобы ожидаемые результаты с выходными как то правильно работали

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

func createMatrixByNN(nn [][]float64) (m [][]float64) {
	m = make([][]float64, len(nn))
	for i, l := range nn {
		m[i] = make([]float64, len(l))
	}
	return m
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
	/*
		//nn := createNeurons(2, 3)
		//m := createMatrixByNN(n1n)
		//fmt.Printf("%v\n\n%v", nn, m)
		//nn := make([][]float64, 2)
		//nn[0] = []float64{0.5, 0.2}
		//nn[1] = []float64{0, 0, 0}
		//
		//
		//fmt.Printf("NN : %v\n", nn)
		//fmt.Printf("W : %v\n", w)
		//forward(nn[0], nn[1], w[0])
		//fmt.Printf("NN : %v\n", nn)
		//fmt.Printf("W : %v\n", w)
	*/
	neurons := createNeurons(2, 3)
	w := createWeights(neurons)
	w[0][0][0] = 1
	w[0][0][1] = 2
	w[0][0][2] = 3
	w[0][1][0] = 4
	w[0][1][1] = 5
	w[0][1][2] = 6
	var nn = NeuralNetwork{neurons: neurons, w: w, LR: 0.1, EPOCH: 50}
	nn.neurons[0] = []float64{0.3, 0.1}

	//forward(nn)
	predict(nn, []float64{1, 0.2})
	fmt.Printf("NN:\n %v\n", nn.neurons)
	fmt.Printf("W:\n %v\n", nn.w)

	nn.neurons[0] = []float64{1, 0}
	forward(nn)

	fmt.Printf("NN:\n %v\n", nn.neurons)
	fmt.Printf("W:\n %v\n", nn.w)
	/*
		//nn.w[0] = [][]float64{{1, 2, 1}, {1, 2, 0}}
		//nn.w[1] = [][]float64{{1, 2, 5}, {1, 2, 0, 0.5}, {1, 1, 1, 0.5}, {1, 0.2, 1, 0.2}}
		//nn.w[2] = [][]float64{{1, 2, 1, 0.5}, {1, 2, 0, 0.5}, {1, 1, 1, 0.5}, {1, 0.2, 1, 0.2}}
		//nn.w[4] = [][]float64{{1, 2, 1, 0.5}, {1, 2, 0, 0.5}, {1, 1, 1, 0.5}, {1, 0.2, 1, 0.2}}
		//nn.w[5] = [][]float64{{1, 2, 1, 0.5}, {1, 2, 0, 0.5}, {1, 1, 1, 0.5}, {1, 0.2, 1, 0.2}}
		//forward(nn.neurons[0], nn.neurons[1], nn.w[0])
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
	*/
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
}

/*
func forward(inputNeurons []float64, outputNeurons []float64, w [][]float64) { // функция прямого распространения

	for i, _ := range outputNeurons {
		var res float64 = 0
		for j, input := range inputNeurons {
			res += input * w[j][i] // i - к output нейрону, j - из input нейрона
		}
		outputNeurons[i] = activate(res)
	}
}
*/

func forward(nn NeuralNetwork) { // функция прямого распространения
	for n := 0; n < len(nn.w); n++ { // цикл по слоям нейронов
		//forward(nn.neurons[n], nn.neurons[n+1], nn.w[n]) // текущий слой, следующий слой, веса между этими слоями
		for i, _ := range nn.neurons[n+1] {
			var res float64 = 0
			for j, input := range nn.neurons[n] {
				res += input * nn.w[n][j][i] // i - к output нейрону, j - из input нейрона
			}
			nn.neurons[n+1][i] = activate(res)
		}

	}

}

func backProp(nn NeuralNetwork, exp []float64) { // Функция обратного распространения ошибки
	// должен тут вызывать forward() и что-то получать и дальше считать цену ошибки
	forward(nn)

	m := createMatrixByNN(nn.neurons)
	lastLayer := len(nn.neurons) - 1
	for i, n := range nn.neurons[lastLayer] {
		m[lastLayer][i] = (1 / (math.Pow(math.Cosh(n), 2))) * (n - exp[i]) // n - полученное значение, exp - ожидаемое значение
	}

	for i := len(nn.neurons) - 2; i >= 0; i-- { // neuron layer
		for j := 0; j < len(nn.neurons[i]); j++ { // some neuron
			var sum float64
			for k, elem := range nn.neurons[i+1] { // recursive sum for this neuron
				sum += elem * nn.w[i][j][k]
			}
			m[i][j] = sum * (1 / (math.Pow(math.Cosh(nn.neurons[i][j]), 2))) // нашли б

			for k, _ := range nn.neurons[i+1] { // коррекция весов
				deltaW := -(nn.LR * m[i+1][j] * nn.neurons[i][j]) // dW - коррекция весов
				nn.w[i][j][k] = nn.w[i][j][k] - deltaW            // изменяем веса
			}
		}
	}
}

func train(nn NeuralNetwork, data [][]float64, exp [][]float64) {

	for e := 0; e < nn.EPOCH; e++ { // цикл по эпохам

		for d := 0; d < len(data); d++ { // цикл по дата сету // вынести в функцию
			for i := 0; i < len(data); i++ {
				nn.neurons[0][i] = data[d][i] // присваиваем входным нейронам данные из дата-сета
			}
			//forward(nn)
			//
			//for n := 0; n < len(nn.w); n++ { // цикл по слоям весов нейронов
			//	forward(nn.neurons[n], nn.neurons[n+1], nn.w[n]) // текущий слой, следующий слой, веса между этими слоями
			//}
			backProp(nn, exp[d]) // вычисляем ошибку и корректируем веса
		}

		accuracy := evaluate()
		fmt.Printf("Epoch: %d, Accuracy: %.2f\n", e, accuracy)

	}
}

func predict(nn NeuralNetwork, testData []float64) []float64 { // вычислить
	for i, _ := range nn.neurons[0] {
		nn.neurons[0][i] = testData[i]
	}
	forward(nn)

	ll := len(nn.neurons) - 1 // last layer
	outputs := make([]float64, len(nn.neurons[ll]))
	for i, _ := range nn.neurons[ll] {
		outputs[i] = nn.neurons[ll][i]
	}
	return outputs
}

func evaluate() (accuracy float64) {

	return accuracy
}
