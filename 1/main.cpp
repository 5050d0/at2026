#include <iostream>
#include <regex>

using namespace std;

int main() {
    string s = "Заказ №12345 доставлен 15.03.2025, второй заказ №6789 от 27.08.2024";
    std::regex word("(?<=№)\\d+", std::regex_constants::ECMAScript);
    auto word_begin = sregex_iterator(s.begin(), s.end(), word);
    auto word_end = sregex_iterator();
    for (auto it = word_begin; it != word_end; ++it) {
        auto m = *it;
        cout << m.str() << std::endl;
    }
}

