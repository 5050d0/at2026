#include <gtest/gtest.h>

#include "recognizer/Loader.h"
TEST(RecognizerTest, EverythingTest) {
    Loader loader;

    std::vector<std::string> const valid_lines = {
        "int ab= 444;",
        "int Dsasf323 = ab /     5;",
        "long vA2r23=5%6;",
        "int ab = ab % ab  ;",
        "long newvar = ab;",
        "short verynewvar = 5 % ab;"
    };
    std::vector<std::string> const error_lines = {
        "short ab = 1;",
        "long ab = 4;"

    };
    std::vector<std::string> const invalid_lines = {
        "long var1=NotAVar   ;",
        "short va2 = 5 + NotAVar;",
        "a ad=5;",
        "int 3d = ab+Dsasf323;",
        "int a =;",
        "int verynewvar2 = 5 % notawar3;"
    };
    for (auto const &str: valid_lines) {
        for (auto &[fst, snd]: loader.get_recognizers()) {
            EXPECT_EQ(snd.Validate(str), std::make_pair(true, "")) << "line: " << str;
        }
    }
    for (auto const &str: error_lines) {
        for (auto &[fst, snd]: loader.get_recognizers()) {
            auto result = snd.Validate(str);
            EXPECT_EQ(result.first, true) << "line: " << str;
            EXPECT_NE(result.second, "") << "line: " << str;
        }
    }
    for (auto const &str: invalid_lines) {
        for (auto &[fst, snd]: loader.get_recognizers()) {
            EXPECT_EQ(snd.Validate(str), std::make_pair(false, "")) << "line: " << str;
        }
    }
    for (auto &[fst, snd]: loader.get_recognizers()) {
        EXPECT_EQ(snd.Validate("int valid = 5;"), std::make_pair(true, ""));
        EXPECT_EQ(snd.Validate("short valid = 5;").second.empty(), false);
        snd.reset();
        EXPECT_EQ(snd.Validate("short valid = 5;").second.empty(), true);
    }
}
