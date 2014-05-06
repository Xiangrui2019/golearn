package linear_models

/*
#cgo LDFLAGS: -L. -llinear
#include "liblinear_src/linear.h"
*/
import "C"

type Problem struct {
	c_prob C.struct_problem
}

type Parameter struct {
	c_param C.struct_parameter
}

type Model struct {
	c_model *C.struct_model
}

const (
	L2R_LR              = C.L2R_LR
	L2R_L2LOSS_SVC_DUAL = C.L2R_L2LOSS_SVC_DUAL
	L2R_L2LOSS_SVC      = C.L2R_L2LOSS_SVC
	L2R_L1LOSS_SVC_DUAL = C.L2R_L1LOSS_SVC_DUAL
	MCSVM_CS            = C.MCSVM_CS
	L1R_L2LOSS_SVC      = C.L1R_L2LOSS_SVC
	L1R_LR              = C.L1R_LR
	L2R_LR_DUAL         = C.L2R_LR_DUAL
	L2R_L2LOSS_SVR      = C.L2R_L2LOSS_SVR
	L2R_L2LOSS_SVR_DUAL = C.L2R_L2LOSS_SVR_DUAL
	L2R_L1LOSS_SVR_DUAL = C.L2R_L1LOSS_SVR_DUAL
)

func NewParameter(solver_type int, C float64, eps float64) *Parameter {
	param := Parameter{}
	param.c_param.solver_type = C.int(solver_type)
	param.c_param.eps = C.double(eps)
	param.c_param.C = C.double(C)
	param.c_param.nr_weight = C.int(0)
	param.c_param.weight_label = nil
	param.c_param.weight = nil
	param.c_param.p = C.double(0.1)

	return &param
}

func NewProblem(X [][]float64, y []float64, bias float64) *Problem {
	prob := Problem{}
	prob.c_prob.l = C.int(len(X))
	prob.c_prob.n = C.int(len(X[0]) + 1)

	prob.c_prob.x = convert_features(X, bias)
	c_y := make([]C.double, len(y))
	for i := 0; i < len(y); i += 1 {
		c_y[i] = C.double(y[i])
	}
	prob.c_prob.y = &c_y[0]
	prob.c_prob.bias = C.double(-1)

	return &prob
}

func Train(prob *Problem, param *Parameter) *Model {
	return &Model{C.train(&prob.c_prob, &param.c_param)}
}

func Predict(model *Model, x []float64) float64 {
	c_x := convert_vector(x, 0)
	c_y := C.predict(model.c_model, c_x)
	y := float64(c_y)
	return y
}
func convert_vector(x []float64, bias float64) *C.struct_feature_node {
	n_ele := 0
	for i := 0; i < len(x); i += 1 {
		if x[i] > 0 {
			n_ele += 1
		}
	}
	n_ele += 2

	c_x := make([]C.struct_feature_node, n_ele)
	j := 0
	for i := 0; i < len(x); i += 1 {
		if x[i] > 0 {
			c_x[j].index = C.int(i + 1)
			c_x[j].value = C.double(x[i])
			j += 1
		}
	}
	if bias > 0 {
		c_x[j].index = C.int(0)
		c_x[j].value = C.double(0)
		j += 1
	}
	c_x[j].index = C.int(-1)
	return &c_x[0]
}
func convert_features(X [][]float64, bias float64) **C.struct_feature_node {
	n_samples := len(X)
	n_elements := 0

	for i := 0; i < n_samples; i += 1 {
		for j := 0; j < len(X[i]); j += 1 {
			if X[i][j] != 0.0 {
				n_elements += 1
			}
			n_elements += 1 //for bias
		}
	}

	x_space := make([]C.struct_feature_node, n_elements+n_samples)

	cursor := 0
	x := make([]*C.struct_feature_node, n_samples)
	var c_x **C.struct_feature_node

	for i := 0; i < n_samples; i += 1 {
		x[i] = &x_space[cursor]

		for j := 0; j < len(X[i]); j += 1 {
			if X[i][j] != 0.0 {
				x_space[cursor].index = C.int(j + 1)
				x_space[cursor].value = C.double(X[i][j])
				cursor += 1
			}
			if bias > 0 {
				x_space[cursor].index = C.int(0)
				x_space[cursor].value = C.double(bias)
				cursor += 1
			}
		}
		x_space[cursor].index = C.int(-1)
		cursor += 1
	}
	c_x = &x[0]
	return c_x
}
