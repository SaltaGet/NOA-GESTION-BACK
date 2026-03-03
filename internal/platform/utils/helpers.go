package utils

import "strings"

func Ternary[T any](cond bool, a, b T) T {
    if cond {
        return a
    }
    return b
}

    func SplitStrings(ptr *string) []string {
        if ptr == nil || *ptr == "" {
            return []string{}
        }

        return strings.Split(*ptr, ",")
    }