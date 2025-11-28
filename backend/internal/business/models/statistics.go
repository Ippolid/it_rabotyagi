package models

import "time"

// UserStatistics представляет общую статистику пользователя
type UserStatistics struct {
	CoursesEnrolled    int
	CoursesCompleted   int
	OverallProgressPct float32
	Questions          QuestionStatisticsSummary
}

// QuestionStatisticsSummary представляет сводку статистики по вопросам
type QuestionStatisticsSummary struct {
	TotalSolved       int
	CorrectAnswersPct float32
	TotalAttempts     int
	AvgTimeSpent      float32
}

// CourseStatisticsItem представляет статистику по одному курсу
type CourseStatisticsItem struct {
	CourseID         int
	CourseTitle      string
	TotalModules     int
	CompletedModules int
	ProgressPct      float32
	StartedAt        time.Time
	TimeSpent        float32
}

// QuestionStatisticsByCourse представляет статистику по вопросам для одного курса
type QuestionStatisticsByCourse struct {
	CourseID    int
	CourseTitle string
	Statistics  QuestionStatisticsSummary
}
