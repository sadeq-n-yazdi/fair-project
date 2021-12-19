
import copy


def readProjectsFromFile(fileName="projects.txt"):
    with open(fileName) as file:
        theProjects = file.readlines()
        theProjects = [project.strip() for project in theProjects]
    return theProjects


def readStudentPrefrencesFromFile(fileName="students.txt"):
    with open(fileName) as file:
        lines = file.readlines()
        theStudents = [
            {
                "Name": (student.split(",")[0]).strip(),
                "Prefrences":[pref.strip() for pref in student.split(",")[1:]],
                "assigned": False,
            } for student in lines
        ]
    return theStudents


Projects = readProjectsFromFile("projects.csv")
Students = readStudentPrefrencesFromFile("students.csv")


def getProjects():
    return copy.deepcopy(Projects)


def getStudents():
    return copy.deepcopy(Students)
