#include <gtest/gtest.h>

#include "recognizer/Loader.h"
TEST(RecognizerTest, EverythingTest) {
    Loader loader;

    std::vector<std::string> const valid_lines = {
        "int ab= 444;",
        "short Dsasf323 = ab /     5;",
        "long vA2r23=5%6;",
        "int ab = ab % ab  ;"
    };
    std::vector<std::string> const error_lines = {
        "short ab = 1;"
    };
    std::vector<std::string> const invalid_lines = {
        "long var1=NotAVar   ;",
        "short va2 = 5 + NotAVar;",
        "a ad=5;",
        "int 3d = ab+Dsasf323;",
        "int a =;"
    };
    for (auto const &str: valid_lines) {
        for (auto &recognizer: loader.get_recognizers()) {
            EXPECT_EQ(recognizer.second->Recognize(str), std::make_pair(true, "")) << "line: " << str;
        }
    }
    for (auto const &str: error_lines) {
        for (auto &recognizer: loader.get_recognizers()) {
            auto result = recognizer.second->Recognize(str);
            EXPECT_EQ(result.first, false) << "line: " << str;
            EXPECT_NE(result.second, "") << "line: " << str;
        }
    }
    for (auto const &str: invalid_lines) {
        for (auto &recognizer: loader.get_recognizers()) {
            EXPECT_EQ(recognizer.second->Recognize(str), std::make_pair(false, "")) << "line: " << str;
        }
    }
}
