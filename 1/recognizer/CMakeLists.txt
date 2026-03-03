

add_library(recognizer Loader.cpp Loader.h)
target_link_libraries(recognizer Validator)

add_library(Validator Validator.cpp Validator.h)
target_link_libraries(Validator RegexRecognizer)


add_library(IRecognizer INTERFACE IRecognizer.h)
target_include_directories(IRecognizer INTERFACE ${CMAKE_SOURCE_DIR})

#add_subdirectory(FlexRecognizer)
add_subdirectory(RegexRecognizer)
#add_subdirectory(SMCRecognizer)
add_subdirectory(Tests)