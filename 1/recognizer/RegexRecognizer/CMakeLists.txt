
add_library(RegexRecognizer RegexRecognizer.cpp RegexRecognizer.h)

target_link_libraries(RegexRecognizer PUBLIC IRecognizer)
target_include_directories(RegexRecognizer PUBLIC ${CMAKE_SOURCE_DIR})


find_package(Boost REQUIRED COMPONENTS regex)
target_link_libraries(RegexRecognizer PRIVATE Boost::regex)