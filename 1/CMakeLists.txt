cmake_minimum_required(VERSION 3.31)
project(lab1)

set(CMAKE_CXX_STANDARD 23)
include_subdirectory(recognizer)
# set(CMAKE_BUILD_TYPE Debug)
add_executable(lab1 main.cpp)
