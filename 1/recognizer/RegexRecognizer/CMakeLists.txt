
add_library(RegexRecognizer RegexRecognizer.cpp RegexRecognizer.h)

target_link_libraries(RegexRecognizer IRecognizer)
target_include_directories(RegexRecognizer PUBLIC ${CMAKE_SOURCE_DIR})