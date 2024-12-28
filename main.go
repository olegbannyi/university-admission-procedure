package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

const filename = "applicants.txt"

var N int
var M int

type University struct {
	enrolledPerDepartment int
	deparments            []Department
}

type Department struct {
	name       string
	applicants []Applicant
}

type Applicant struct {
	firstName    string
	lastName     string
	exams        []float64
	specialScore float64
	score        float64
	departments  []string
}

func main() {
	fmt.Scan(&N)

	university := NewUniversity(N)
	university.enroll(fetchApplicants())

	for _, dep := range university.deparments {
		f, err := os.Create(strings.ToLower(dep.name) + ".txt")
		if err != nil {
			log.Fatalln(err)
		}

		for _, a := range dep.applicants {
			fmt.Fprintf(f, "%s %.2f\n", a.fullName(), a.score)
		}

		f.Close()
	}
}

func fetchApplicants() map[string]Applicant {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	applicants := make(map[string]Applicant)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")
		exams := make([]float64, 0, 4)

		for i := 0; i < 4; i++ {
			exam, err := strconv.ParseFloat(words[2+i], 64)
			if err != nil {
				log.Fatalln(err)
			}
			exams = append(exams, exam)
		}

		specialScore, err := strconv.ParseFloat(words[6], 64)
		if err != nil {
			log.Fatalln(err)
		}

		applicant := NewApplicant(words[0], words[1], exams, specialScore, words[7], words[8], words[9])
		applicants[applicant.fullName()] = applicant
	}

	return applicants
}

func NewApplicant(firstName string,
	lastName string,
	exams []float64,
	specialScore float64,
	fitstDepartment string,
	secondDepartment string,
	thirdDepartment string) Applicant {
	applicant := Applicant{firstName: firstName,
		lastName:     lastName,
		exams:        exams,
		specialScore: specialScore,
		departments:  []string{fitstDepartment, secondDepartment, thirdDepartment},
	}

	return applicant
}

func (a Applicant) fullName() string {
	return a.firstName + " " + a.lastName
}

func (a Applicant) scoreByDepartment(department string) float64 {
	score := 0.0
	exans := mainExamIndex(department)
	for _, i := range exans {
		score += a.exams[i]
	}
	return math.Max(score/float64(len(exans)), a.specialScore)
}

func NewUniversity(enrolledPerDepartment int) *University {
	univesity := University{enrolledPerDepartment: enrolledPerDepartment, deparments: make([]Department, 5)}

	univesity.deparments[0] = NewDepartment("Biotech", enrolledPerDepartment)
	univesity.deparments[1] = NewDepartment("Chemistry", enrolledPerDepartment)
	univesity.deparments[2] = NewDepartment("Engineering", enrolledPerDepartment)
	univesity.deparments[3] = NewDepartment("Mathematics", enrolledPerDepartment)
	univesity.deparments[4] = NewDepartment("Physics", enrolledPerDepartment)

	return &univesity
}

func NewDepartment(name string, numberOfEnrolled int) Department {
	return Department{name: name, applicants: make([]Applicant, 0, numberOfEnrolled)}
}

func (u *University) enroll(applicants map[string]Applicant) {
	for i := 0; i < 3; i++ {
		for n, dep := range u.deparments {
			if len(u.deparments[n].applicants) < N {
				rest := N - len(dep.applicants)
				depApplicants := filterApplicantsByDepartment(applicants, dep.name, i)

				for j := 0; j < len(depApplicants) && j < rest; j++ {
					depApplicants[j].score = depApplicants[j].scoreByDepartment(dep.name)
					u.deparments[n].applicants = append(u.deparments[n].applicants, depApplicants[j])
					delete(applicants, depApplicants[j].fullName())
				}
			}
		}
	}

	for n := range u.deparments {
		sort.Slice(u.deparments[n].applicants, func(i, j int) bool {
			if u.deparments[n].applicants[i].score != u.deparments[n].applicants[j].score {
				return u.deparments[n].applicants[i].score > u.deparments[n].applicants[j].score
			}
			return u.deparments[n].applicants[i].fullName() < u.deparments[n].applicants[j].fullName()
		})
	}
}

func departmentIndex(a Applicant, department string) int {
	for i, d := range a.departments {
		if d == department {
			return i
		}
	}
	return -1
}

func filterApplicantsByDepartment(applicants map[string]Applicant, department string, level int) []Applicant {
	filtered := make([]Applicant, 0)

	for _, a := range applicants {
		if a.departments[level] == department {
			filtered = append(filtered, a)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].scoreByDepartment(department) == filtered[j].scoreByDepartment(department) {
			return filtered[i].fullName() < filtered[j].fullName()
		}
		return filtered[i].scoreByDepartment(department) > filtered[j].scoreByDepartment(department)
	})

	return filtered
}

func mainExamIndex(departmentName string) []int {
	switch departmentName {
	case "Engineering":
		return []int{2, 3}
	case "Mathematics":
		return []int{2}
	case "Physics":
		return []int{0, 2}
	case "Chemistry":
		return []int{1}
	case "Biotech":
		return []int{0, 1}
	default:
		return []int{}
	}
}
