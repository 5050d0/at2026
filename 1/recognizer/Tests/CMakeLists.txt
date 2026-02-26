enable_testing()

add_executable(
        Recognizer_test
        Recognizer_test.cc
)
target_link_libraries(
        Recognizer_test PUBLIC
        GTest::gtest_main
        recognizer
)

include(GoogleTest)
gtest_discover_tests(Recognizer_test)
