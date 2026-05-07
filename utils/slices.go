package utils

func MapSliceErr[A, B any](slice []A, f func(A) (B, error)) ([]B, error) {
	result := make([]B, 0, len(slice))
	for _, a := range slice {
		b, err := f(a)
		if err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}
