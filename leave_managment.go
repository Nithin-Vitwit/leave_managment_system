package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

const dateFormat = "02-01-2006"
const hrPassword = "admin"
const dataFileName = "data.json"

type System struct {
	Employees map[string]*Emp `json:"employees"`
	HRs       map[string]*Hr  `json:"hrs"`
	Leaves    []*Leave        `json:"leaves"`
	DataFile  string          `json:"-"`
}

type member interface {
	showActions(s *System)
}

type Leave struct {
	Name     string    `json:"name"`
	EmpID    string    `json:"emp_id"`
	Reason   string    `json:"reason"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	Status   string    `json:"status"`
}

type Emp struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Hr struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func loadSystem(fileName string) *System {
	s := &System{
		Employees: make(map[string]*Emp),
		HRs:       make(map[string]*Hr),
		Leaves:    make([]*Leave, 0),
		DataFile:  fileName,
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf(" %s not found. Starting with a fresh system.\n", fileName)
			return s
		}
		fmt.Printf("Error reading file %s: %v. Starting fresh.\n", fileName, err)
		return s
	}

	if err := json.Unmarshal(data, s); err != nil {
		fmt.Printf("Error unmarshalling JSON from %s: %v. Starting fresh.\n", fileName, err)
		return s
	}

	fmt.Printf("System state loaded from %s.\n", fileName)
	return s
}

func saveSystem(s *System) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling system data: %v\n", err)
		return
	}

	if err := os.WriteFile(s.DataFile, data, 0644); err != nil {
		fmt.Printf("Error writing system data to %s: %v\n", s.DataFile, err)
	} else {
		fmt.Printf("💾 System state saved to %s.\n", s.DataFile)
	}
}

func (e *Emp) showActions(s *System) {
	var action int
	fmt.Println("\n--- Employee Menu ---")
	fmt.Printf("Welcome, %s (%s)\n", e.Name, e.ID)
	fmt.Println("Enter the number for the action to perform:")
	fmt.Println("1) Apply for leave")
	fmt.Println("2) View my leave history")
	fmt.Println("3) Switch user/Exit")

	fmt.Scan(&action)
	switch action {
	case 1:
		e.applyLeave(s)
	case 2:
		e.viewLeaveHistory(s)
	case 3:
		login(s)
	default:
		fmt.Println("Invalid action.")
		e.showActions(s)
	}
}

func (e *Emp) applyLeave(s *System) {
	e.newLeave(s)
}

func (e *Emp) newLeave(s *System) {
	var reason, fromDateStr, toDateStr string

	fmt.Println("--- Apply for Leave ---")
	fmt.Print("Give the reason for the leave: ")
	if _, err := fmt.Scanln(&reason); err != nil {
	}

	fmt.Printf("Give **from date** (format: %s): ", dateFormat)
	fmt.Scan(&fromDateStr)
	fromDate, err := time.Parse(dateFormat, fromDateStr)
	if err != nil {
		fmt.Println("Error parsing From Date. Please use the correct format.")
		e.showActions(s)
		return
	}

	fmt.Printf("Give **to date** (format: %s): ", dateFormat)
	fmt.Scan(&toDateStr)
	toDate, err := time.Parse(dateFormat, toDateStr)
	if err != nil {
		fmt.Println("Error parsing To Date. Please use the correct format.")
		e.showActions(s)
		return
	}

	l := &Leave{
		Name:     e.Name,
		EmpID:    e.ID,
		Reason:   reason,
		FromDate: fromDate,
		ToDate:   toDate,
		Status:   "Pending",
	}

	s.Leaves = append(s.Leaves, l)
	saveSystem(s)
	fmt.Println("\nLeave application submitted successfully! Status: Pending")
	e.showActions(s)
}

func (e *Emp) viewLeaveHistory(s *System) {
	fmt.Println("\n--- My Leave History ---")
	found := false
	for _, l := range s.Leaves {
		if l.EmpID == e.ID {
			fmt.Printf("Reason: %s | Dates: %s to %s | Status: **%s**\n",
				l.Reason, l.FromDate.Format(dateFormat), l.ToDate.Format(dateFormat), l.Status)
			found = true
		}
	}
	if !found {
		fmt.Println("No leave applications found.")
	}
	e.showActions(s)
}

func (h *Hr) showActions(s *System) {
	hrLogin(s)
}

func hrLogin(s *System) {
	var password string
	fmt.Println("\n--- HR Login ---")
	fmt.Print("Enter the password: ")
	fmt.Scan(&password)

	if password == hrPassword {
		if _, ok := s.HRs["12"]; ok {
			h := s.HRs["12"]
			h.showPendingLeaves(s)
		} else {
			fmt.Println("Error: HR user (ID 12) not found in system data.")
			login(s)
		}
	} else {
		fmt.Println(" Invalid password. Returning to main login.")
		login(s)
	}
}

func (h *Hr) showPendingLeaves(s *System) {
	fmt.Println("\n--- Pending Leave Applications ---")
	pendingLeavesMap := make(map[int]*Leave)
	displayCount := 1

	for i, l := range s.Leaves {
		if l.Status == "Pending" {
			fmt.Printf("%d) **Name:** %s (ID: %s) | **Reason:** %s | **Dates:** %s to %s | **System Index:** %d\n",
				displayCount, l.Name, l.EmpID, l.Reason, l.FromDate.Format(dateFormat), l.ToDate.Format(dateFormat), i)
			pendingLeavesMap[displayCount] = l
			displayCount++
		}
	}

	if displayCount == 1 {
		fmt.Println("🎉 No pending leave applications.")
		login(s)
		return
	}

	var selectionStr string
	fmt.Print("\nTo grant or decline, select the **number** (1 to %d) of the leave (or type 'exit'): ", displayCount-1)
	fmt.Scan(&selectionStr)

	if selectionStr == "exit" {
		login(s)
		return
	}

	selectionIndex, err := strconv.Atoi(selectionStr)
	selectedLeave := pendingLeavesMap[selectionIndex]

	if err != nil || selectedLeave == nil {
		fmt.Println("Invalid selection. Please enter a valid number for a pending leave.")
		h.showPendingLeaves(s)
		return
	}

	var decision int
	fmt.Printf("Selected leave for %s. Select **1** for decline or **2** for accepting: ", selectedLeave.Name)
	fmt.Scan(&decision)

	if decision == 1 {
		h.decline(selectedLeave, s)
	} else if decision == 2 {
		h.grant(selectedLeave, s)
	} else {
		fmt.Println("Invalid decision. Please try again.")
	}

	h.showPendingLeaves(s)
}

func (h *Hr) grant(l *Leave, s *System) {
	l.Status = "Granted"
	saveSystem(s)
	fmt.Printf("✅ Leave for %s has been **GRANTED**.\n", l.Name)
}

func (h *Hr) decline(l *Leave, s *System) {
	l.Status = "Declined"
	saveSystem(s)
	fmt.Printf("❌ Leave for %s has been **DECLINED**.\n", l.Name)
}

func login(s *System) {
	var userID string
	var choice int

	fmt.Println("\n==================================")
	fmt.Println("💼 Leave Management System Login 📅")
	fmt.Println("==================================")
	fmt.Println("Select User Type:")
	fmt.Println("1) Employee")
	fmt.Println("2) HR")
	fmt.Println("3) Exit Program")
	fmt.Print("Enter choice: ")

	if _, err := fmt.Scan(&choice); err != nil {
		fmt.Println("Error reading input. Exiting.")
		os.Exit(1)
	}

	switch choice {
	case 1:
		fmt.Print("Enter Employee ID: ")
		fmt.Scan(&userID)
		if emp, ok := s.Employees[userID]; ok {
			emp.showActions(s)
		} else {
			fmt.Println(" Employee ID not found.")
			login(s)
		}
	case 2:
		hrLogin(s)
	case 3:
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice.")
		login(s)
	}
}

func main() {
	system := loadSystem(dataFileName)

	if len(system.Employees) == 0 {
		fmt.Println("Initializing default users...")
		emp1 := &Emp{Name: "Nithin", ID: "22"}
		emp2 := &Emp{Name: "Kalyan", ID: "23"}
		hr1 := &Hr{Name: "Soumya", ID: "12"}

		system.Employees["22"] = emp1
		system.Employees["23"] = emp2
		system.HRs["12"] = hr1

		saveSystem(system)
	}

	login(system)
}
