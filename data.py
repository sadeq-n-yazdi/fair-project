
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


Projects = readProjectsFromFile()
Students = readStudentPrefrencesFromFile()


# Students = [
#     {
#         "Name": "Student 1",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[0],
#             Projects[3],
#             Projects[1],
#         ]
#     },
#     {
#         "Name": "Student 2",
#         "Prefrences": [
#             Projects[1],
#             Projects[2],
#             Projects[3],
#             Projects[5],
#             Projects[4],
#         ]
#     },
#     {
#         "Name": "Student 3",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[0],
#             Projects[4],
#             Projects[1],
#         ]
#     },
#     {
#         "Name": "Student 4",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[0],
#             Projects[3],
#             Projects[1],
#         ]
#     },
#     {
#         "Name": "Student 5",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[4],
#             Projects[29],
#             Projects[1],
#         ]
#     },
#     {
#         "Name": "Student 6",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[0],
#             Projects[3],
#             Projects[11],
#         ]
#     },

#     {
#         "Name": "Student 7",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[0],
#             Projects[15],
#             Projects[1],
#         ]
#     },

#     {
#         "Name": "Student 8",
#         "Prefrences": [
#             Projects[6],
#             Projects[2],
#             Projects[10],
#             Projects[3],
#             Projects[12],
#         ]
#     },

#     {
#         "Name": "Student 9",
#         "Prefrences": [
#             Projects[5],
#             Projects[21],
#             Projects[0],
#             Projects[13],
#             Projects[1],
#         ]
#     },

#     {
#         "Name": "Student 10",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[0],
#             Projects[3],
#             Projects[1],
#         ]
#     },

#     {
#         "Name": "Student 11",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[7],
#             Projects[3],
#             Projects[1],
#         ]
#     },

#     {
#         "Name": "Student 12",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[7],
#             Projects[8],
#             Projects[9],
#         ]
#     },

#     {
#         "Name": "Student 13",
#         "Prefrences": [
#             Projects[5],
#             Projects[2],
#             Projects[7],
#             Projects[8],
#             Projects[9],
#         ]
#     },

# ]

def getProjects():
    return copy.deepcopy(Projects)


def getStudents():
    return copy.deepcopy(Students)
