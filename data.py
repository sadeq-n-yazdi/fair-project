
import copy
Projects = [
    "Proj 01",
    "Proj 02",
    "Proj 03",
    "Proj 04",
    "Proj 05",
    "Proj 06",
    "Proj 07",
    "Proj 08",
    "Proj 09",
    "Proj 10",
    "Proj 11",
    "Proj 12",
    "Proj 13",
    "Proj 14",
    "Proj 15",
    "Proj 16",
    "Proj 17",
    "Proj 18",
    "Proj 19",
    "Proj 20",
    "Proj 21",
    "Proj 22",
    "Proj 23",
    "Proj 24",
    "Proj 25",
    "Proj 26",
    "Proj 27",
    "Proj 28",
    "Proj 29",
    "Proj 30",
]

Students = [
    {
        "Name": "Student 1",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[0],
            Projects[3],
            Projects[1],
        ]
    },
    {
        "Name": "Student 2",
        "Prefrences": [
            Projects[1],
            Projects[2],
            Projects[3],
            Projects[5],
            Projects[4],
        ]
    },
    {
        "Name": "Student 3",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[0],
            Projects[4],
            Projects[1],
        ]
    },
    {
        "Name": "Student 4",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[0],
            Projects[3],
            Projects[1],
        ]
    },
    {
        "Name": "Student 5",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[4],
            Projects[29],
            Projects[1],
        ]
    },
    {
        "Name": "Student 6",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[0],
            Projects[3],
            Projects[11],
        ]
    },

    {
        "Name": "Student 7",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[0],
            Projects[15],
            Projects[1],
        ]
    },

    {
        "Name": "Student 8",
        "Prefrences": [
            Projects[6],
            Projects[2],
            Projects[10],
            Projects[3],
            Projects[12],
        ]
    },

    {
        "Name": "Student 9",
        "Prefrences": [
            Projects[5],
            Projects[21],
            Projects[0],
            Projects[13],
            Projects[1],
        ]
    },

    {
        "Name": "Student 10",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[0],
            Projects[3],
            Projects[1],
        ]
    },

    {
        "Name": "Student 11",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[7],
            Projects[3],
            Projects[1],
        ]
    },

    {
        "Name": "Student 12",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[7],
            Projects[8],
            Projects[9],
        ]
    },

    {
        "Name": "Student 13",
        "Prefrences": [
            Projects[5],
            Projects[2],
            Projects[7],
            Projects[8],
            Projects[9],
        ]
    },

]

# print (Students)


def getProjects():
    return copy.deepcopy(Projects)


def getStudents():
    return copy.deepcopy(Students)
