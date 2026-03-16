

add_library(recognizer Loader.cpp Loader.h)
add_library(Validator Validator.cpp Validator.h)
add_library(IRecognizer INTERFACE IRecognizer.h)

target_include_directories(IRecognizer INTERFACE ${CMAKE_SOURCE_DIR})

target_link_libraries(Validator IRecognizer)
target_link_libraries(recognizer PUBLIC Validator RegexRecognizer SmcRecognizer FlexRecognizer)


add_subdirectory(FlexRecognizer)
add_subdirectory(RegexRecognizer)
add_subdirectory(SMCRecognizer)
add_subdirectory(Tests)