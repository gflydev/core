package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gflydev/core/utils"
)

// This example demonstrates how to use the utils package in a real application context.
// It shows various utility functions for string manipulation, file operations, time handling, and more.

func main() {
	fmt.Println("gFly Core Utils Package Examples")
	fmt.Println("==============================")

	// Example 1: String utilities
	fmt.Println("\nExample 1: String Utilities")
	fmt.Println("------------------------")
	stringUtilitiesExample()

	// Example 2: File utilities
	fmt.Println("\nExample 2: File Utilities")
	fmt.Println("----------------------")
	fileUtilitiesExample()

	// Example 3: Time utilities
	fmt.Println("\nExample 3: Time Utilities")
	fmt.Println("----------------------")
	timeUtilitiesExample()

	// Example 4: Slice utilities
	fmt.Println("\nExample 4: Slice Utilities")
	fmt.Println("-----------------------")
	sliceUtilitiesExample()

	// Example 5: Map utilities
	fmt.Println("\nExample 5: Map Utilities")
	fmt.Println("---------------------")
	mapUtilitiesExample()

	// Example 6: Environment utilities
	fmt.Println("\nExample 6: Environment Utilities")
	fmt.Println("-----------------------------")
	environmentUtilitiesExample()

	// Example 7: Password utilities
	fmt.Println("\nExample 7: Password Utilities")
	fmt.Println("---------------------------")
	passwordUtilitiesExample()
}

func stringUtilitiesExample() {
	// String case conversion
	original := "hello_world_example"
	camelCase := utils.ToCamelCase(original)
	pascalCase := utils.ToPascalCase(original)
	snakeCase := utils.ToSnakeCase("HelloWorldExample")

	fmt.Printf("Original: %s\n", original)
	fmt.Printf("Camel case: %s\n", camelCase)
	fmt.Printf("Pascal case: %s\n", pascalCase)
	fmt.Printf("Snake case: %s\n", snakeCase)

	// String checking
	fmt.Printf("\nIs empty string empty? %t\n", utils.IsEmpty(""))
	fmt.Printf("Is 'hello' empty? %t\n", utils.IsEmpty("hello"))
	fmt.Printf("Is blank string blank? %t\n", utils.IsBlank(""))
	fmt.Printf("Is '   ' blank? %t\n", utils.IsBlank("   "))
	fmt.Printf("Is 'hello' blank? %t\n", utils.IsBlank("hello"))

	// String truncation
	longString := "This is a very long string that needs to be truncated to a shorter length"
	truncated := utils.Truncate(longString, 20, true)
	fmt.Printf("\nOriginal: %s\n", longString)
	fmt.Printf("Truncated: %s\n", truncated)

	// String inclusion in slice
	fruits := []string{"apple", "banana", "orange", "grape"}
	fmt.Printf("\nFruits: %v\n", fruits)
	fmt.Printf("Does fruits include 'banana'? %t\n", utils.IncludeStr(fruits, "banana"))
	fmt.Printf("Does fruits include 'kiwi'? %t\n", utils.IncludeStr(fruits, "kiwi"))
	fmt.Printf("Index of 'orange' in fruits: %d\n", utils.IndexOfStr(fruits, "orange"))
	fmt.Printf("Index of 'kiwi' in fruits: %d\n", utils.IndexOfStr(fruits, "kiwi"))
}

func fileUtilitiesExample() {
	// Create a temporary directory for file operations
	tempDir, err := os.MkdirTemp("", "gfly-utils-example")
	if err != nil {
		fmt.Printf("Failed to create temp directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	content := "Hello, World! This is a test file."
	err = utils.WriteStringToFile(testFilePath, content, 0644)
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
		return
	}

	// File existence and properties
	fmt.Printf("File exists: %t\n", utils.FileExists(testFilePath))
	fmt.Printf("Directory exists: %t\n", utils.DirExists(tempDir))
	fmt.Printf("File extension: %s\n", utils.FileExt(testFilePath))
	fmt.Printf("File size: %d bytes\n", utils.FileSize(testFilePath))

	// Read file content
	fileContent, err := utils.ReadFileAsString(testFilePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
	} else {
		fmt.Printf("File content: %s\n", fileContent)
	}

	// Copy file
	copiedFilePath := filepath.Join(tempDir, "test_copy.txt")
	err = utils.CopyFile(testFilePath, copiedFilePath, 0644)
	if err != nil {
		fmt.Printf("Failed to copy file: %v\n", err)
	} else {
		fmt.Printf("File copied to: %s\n", copiedFilePath)
		fmt.Printf("Copied file exists: %t\n", utils.FileExists(copiedFilePath))
	}

	// Rename file
	newName := utils.RenameFile(testFilePath, "renamed.txt")
	renamedFilePath := filepath.Join(tempDir, newName)
	err = os.Rename(testFilePath, renamedFilePath)
	if err != nil {
		fmt.Printf("Failed to rename file: %v\n", err)
	} else {
		fmt.Printf("File renamed to: %s\n", newName)
		fmt.Printf("Renamed file exists: %t\n", utils.FileExists(renamedFilePath))
	}

	// Create a directory if it doesn't exist
	newDirPath := filepath.Join(tempDir, "new_dir")
	err = utils.MkdirIfNotExists(newDirPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
	} else {
		fmt.Printf("Directory created: %t\n", utils.DirExists(newDirPath))
	}

	// Get file modification time
	modTime, err := utils.FileModTime(copiedFilePath)
	if err != nil {
		fmt.Printf("Failed to get file modification time: %v\n", err)
	} else {
		fmt.Printf("File modification time: %s\n", modTime.Format(time.RFC3339))
	}

	// Check if one file is newer than another
	time.Sleep(1 * time.Second) // Ensure there's a time difference
	err = utils.WriteStringToFile(renamedFilePath, "Updated content", 0644)
	if err != nil {
		fmt.Printf("Failed to update file: %v\n", err)
	}

	isNewer, err := utils.IsFileNewer(renamedFilePath, copiedFilePath)
	if err != nil {
		fmt.Printf("Failed to compare file times: %v\n", err)
	} else {
		fmt.Printf("Is renamed file newer than copied file? %t\n", isNewer)
	}
}

func timeUtilitiesExample() {
	// Current time
	now := utils.Now()
	utc := utils.UTC()
	fmt.Printf("Current local time: %s\n", now.Format(time.RFC3339))
	fmt.Printf("Current UTC time: %s\n", utc.Format(time.RFC3339))

	// Time formatting and parsing
	timeStr := "2023-04-01T14:30:00Z"
	parsedTime, err := utils.ParseTime(timeStr, time.RFC3339)
	if err != nil {
		fmt.Printf("Failed to parse time: %v\n", err)
	} else {
		fmt.Printf("\nParsed time: %s\n", parsedTime.Format(time.RFC3339))
		formatted := utils.FormatTime(parsedTime, "2006-01-02 15:04:05")
		fmt.Printf("Formatted time: %s\n", formatted)
	}

	// Day, week, and month boundaries
	fmt.Printf("\nStart of day: %s\n", utils.StartOfDay(now).Format(time.RFC3339))
	fmt.Printf("End of day: %s\n", utils.EndOfDay(now).Format(time.RFC3339))
	fmt.Printf("Start of week: %s\n", utils.StartOfWeek(now).Format(time.RFC3339))
	fmt.Printf("End of week: %s\n", utils.EndOfWeek(now).Format(time.RFC3339))
	fmt.Printf("Start of month: %s\n", utils.StartOfMonth(now).Format(time.RFC3339))
	fmt.Printf("End of month: %s\n", utils.EndOfMonth(now).Format(time.RFC3339))
	fmt.Printf("Days in month: %d\n", utils.DaysInMonth(now))

	// Day calculations
	tomorrow := utils.AddDays(now, 1)
	yesterday := utils.SubtractDays(now, 1)
	fmt.Printf("\nTomorrow: %s\n", tomorrow.Format("2006-01-02"))
	fmt.Printf("Yesterday: %s\n", yesterday.Format("2006-01-02"))
	fmt.Printf("Days between yesterday and tomorrow: %d\n", utils.DaysBetween(yesterday, tomorrow))

	// Day type checks
	saturday := time.Date(2023, 4, 1, 12, 0, 0, 0, time.Local) // April 1, 2023 was a Saturday
	monday := time.Date(2023, 4, 3, 12, 0, 0, 0, time.Local)   // April 3, 2023 was a Monday
	fmt.Printf("\nIs Saturday a weekend? %t\n", utils.IsWeekend(saturday))
	fmt.Printf("Is Monday a weekend? %t\n", utils.IsWeekend(monday))
	fmt.Printf("Is Monday a weekday? %t\n", utils.IsWeekday(monday))

	// Time comparison
	date1 := time.Date(2023, 4, 1, 12, 0, 0, 0, time.Local)
	date2 := time.Date(2023, 4, 1, 18, 0, 0, 0, time.Local)
	date3 := time.Date(2023, 5, 1, 12, 0, 0, 0, time.Local)
	fmt.Printf("\nAre dates on the same day? %t\n", utils.IsSameDay(date1, date2))
	fmt.Printf("Are dates in the same month? %t\n", utils.IsSameMonth(date1, date3))
	fmt.Printf("Are dates in the same year? %t\n", utils.IsSameYear(date1, date3))
}

func sliceUtilitiesExample() {
	// Create some test slices
	intSlice := []int{1, 2, 3, 4, 5}
	stringSlice := []string{"apple", "banana", "orange", "grape"}

	// Check if a slice contains an element
	fmt.Printf("Does stringSlice contain 'banana'? %t\n", utils.IncludeStr(stringSlice, "banana"))
	fmt.Printf("Does stringSlice contain 'kiwi'? %t\n", utils.IncludeStr(stringSlice, "kiwi"))

	// Find the index of an element in a slice
	fmt.Printf("Index of 'orange' in stringSlice: %d\n", utils.IndexOfStr(stringSlice, "orange"))
	fmt.Printf("Index of 'kiwi' in stringSlice: %d\n", utils.IndexOfStr(stringSlice, "kiwi"))

	// Note: The utils package has more slice utilities that could be demonstrated here
	// if they were available in the package
	fmt.Printf("\nInt slice: %v\n", intSlice)
	fmt.Printf("String slice: %v\n", stringSlice)
}

func mapUtilitiesExample() {
	// Create a test map
	testMap := map[string]interface{}{
		"name":    "John Doe",
		"age":     30,
		"email":   "john@example.com",
		"active":  true,
		"address": map[string]string{"city": "New York", "country": "USA"},
	}

	// Note: The utils package has map utilities that could be demonstrated here
	// if they were available in the package
	fmt.Printf("Test map: %v\n", testMap)

	// For now, just demonstrate basic map operations
	fmt.Printf("\nMap keys: ")
	for key := range testMap {
		fmt.Printf("%s ", key)
	}
	fmt.Println()

	// Check if a key exists
	key := "email"
	if val, ok := testMap[key]; ok {
		fmt.Printf("Value for key '%s': %v\n", key, val)
	} else {
		fmt.Printf("Key '%s' not found\n", key)
	}
}

func environmentUtilitiesExample() {
	// Set some environment variables for testing
	os.Setenv("GFLY_TEST_STRING", "hello world")
	os.Setenv("GFLY_TEST_INT", "42")
	os.Setenv("GFLY_TEST_BOOL", "true")
	os.Setenv("GFLY_TEST_FLOAT", "3.14")

	// Note: The utils package has environment utilities that could be demonstrated here
	// if they were available in the package
	fmt.Println("Environment variables:")
	fmt.Printf("GFLY_TEST_STRING: %s\n", os.Getenv("GFLY_TEST_STRING"))
	fmt.Printf("GFLY_TEST_INT: %s\n", os.Getenv("GFLY_TEST_INT"))
	fmt.Printf("GFLY_TEST_BOOL: %s\n", os.Getenv("GFLY_TEST_BOOL"))
	fmt.Printf("GFLY_TEST_FLOAT: %s\n", os.Getenv("GFLY_TEST_FLOAT"))
}

func passwordUtilitiesExample() {
	// Note: The utils package has password utilities that could be demonstrated here
	// if they were available in the package
	password := "MySecurePassword123!"

	fmt.Printf("Password: %s\n", password)
	fmt.Println("Password utilities would be demonstrated here if available")

	// Example of password validation (simplified)
	isValid := len(password) >= 8 &&
		strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") &&
		strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") &&
		strings.ContainsAny(password, "0123456789") &&
		strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")

	fmt.Printf("Is password valid? %t\n", isValid)
}
