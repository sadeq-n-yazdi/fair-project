#!/bin/python3

import data
import random
import pprint

selectedByProjects = dict()
selectedByStudent = dict()
projects = data.getProjects()
students = data.getStudents()

for s in students:
    s["badPrefrences"] = []
    s["assigned"] = False

for i in range(0, max(len(projects), len(students))):
    random.shuffle(students)

    for student in students:
        if "assigned" in student.keys() and student["assigned"]:
            continue
        if len(student["Prefrences"]) == 0:
            continue

        firstChoice = student["Prefrences"][0]
        studentName = student["Name"]
        if firstChoice in projects:
            selectedByProjects[firstChoice] = studentName
            selectedByStudent[studentName] = firstChoice
            projects.remove(firstChoice)
            student["assigned"] = True
            student["assignedProject"] = firstChoice
        else:
            student["badPrefrences"].append(student["Prefrences"].pop(0))

print()
print("========================")
pprint.pprint(selectedByProjects)
print("------------------------")
pprint.pprint(selectedByStudent)
# print("------------------------")
# print(projects)
# print("------------------------")
# print("------------------------")
# pprint.pprint(students, width=300)
# print("------------------------")
print("------------------------\n")
