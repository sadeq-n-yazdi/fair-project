#!/bin/python3

from data import Projects, Students
import random
import copy
import pprint

selectedByProjects = dict()
selectedByStudent = dict()
projects = copy.deepcopy(Projects)
students = copy.deepcopy(Students)
for s in students:
    s["badPrefrences"] = []
    s["assigned"] = False

for i in range(1, len(Students)):
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
        else:
            student["badPrefrences"].append(student["Prefrences"].pop(0))

print()
print("========================")
print(selectedByProjects)
print("------------------------")
print(selectedByStudent)
print("------------------------")
print(projects)
print("------------------------")
print("------------------------")
pprint.pprint(students, width=300)
print("------------------------")
print("------------------------")
