/*
Copyright (C) 2017 Verizon. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ast

type LogicalExpression struct{ BasicExpression }

func (e *LogicalExpression) Set(input Input) error {
	return e.Initialize(input)
}

func (e *LogicalExpression) Negate() {
	switch e.Operator {
	case "or":
		e.Operator = "and"
	default:
		e.Operator = "or"
	}
	for _, expression := range e.Subexpressions {
		expression.Negate()
	}
}
