package stores

import "github.com/google/uuid"

type Memory struct {
	employees map[string]Employee
}

func NewMemory() *Memory {
	m := Memory{}
	m.employees = map[string]Employee{}
	return &m
}

func (m Memory) Create(_ string, employee *Employee) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	employee.ID = id.String()
	m.employees[employee.ID] = *employee
	return nil
}

func (m Memory) Connect() error {
	return nil
}

func (m Memory) Delete(_ string, id string) error {
	delete(m.employees, id)
	return nil
}

func (m Memory) Update(listID string, id string, newT *Employee) (*Employee, error) {
	oldT, err := m.Get(listID, id)
	if err != nil {
		return nil, err
	}
	if oldT != nil {
		if newT.First_Name != "" {
			oldT.First_Name = newT.First_Name
		}
		if newT.Last_Name != "" {
			oldT.Last_Name = newT.Last_Name
		}
		if newT.Department != "" {
			oldT.Department = newT.Department
		}
		if newT.Salary != 0 {
			oldT.Salary = newT.Salary
		}
		if newT.Age != 0 {
			oldT.Age = newT.Age
		}
		m.employees[id] = *oldT
		return oldT, nil
	}

	return nil, nil
}

func (m Memory) Get(_ string, id string) (*Employee, error) {
	t := m.employees[id]
	return &t, nil
}

func (m Memory) Clear(_ string) error {
	for k := range m.employees {
		delete(m.employees, k)
	}
	return nil
}

func (m Memory) List(_ string) ([]Employee, error) {
	result := []Employee{}
	for _, t := range m.employees {
		result = append(result, t)
	}
	return result, nil
}
