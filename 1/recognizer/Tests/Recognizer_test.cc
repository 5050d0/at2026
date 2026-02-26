#include <gtest/gtest.h>

#include "recognizer/Loader.h"
TEST(RecognizerTest, EverythingTest) {
    Loader loader;

    std::vector<std::string> const valid_lines = {
        "int ab= frgfr;",
        "short Dsasf323 = ab /     5;",
        "long vA2r23=5%6;",
        "int ab = ab % ab  ;"
    };
    for (auto const &str: strings) {
    }
}
