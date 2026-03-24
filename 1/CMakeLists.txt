cmake_minimum_required(VERSION 3.31)
project(lab1)

set(CMAKE_CXX_STANDARD 23)

#set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -g -O0")

include(FetchContent)
FetchContent_Declare(
        googletest
        URL https://github.com/google/googletest/releases/download/v1.17.0/googletest-1.17.0.tar.gz
)
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)


add_subdirectory(recognizer)
# set(CMAKE_BUILD_TYPE Debug)
add_executable(lab1 main.cpp)
target_link_libraries(lab1 recognizer)

add_executable(speedtest speedtest.cpp)
target_link_libraries(speedtest recognizer)