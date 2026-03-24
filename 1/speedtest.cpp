
#include <vector>
#include <iostream>
#include <sstream>
#include <random>

#include "recognizer/Loader.h"

std::string get_vartype(size_t rnd) {
    std::vector<std::string> data = {
        "int", "long", "short"
    };
    return data[rnd % data.size()];
}

std::string get_spaces(size_t length) {
    std::string spaces;
    spaces.reserve(length);
    for (size_t i = 0; i < length; i++) {
        spaces += " ";
    }
    return spaces;
}

size_t get_part(size_t &length, size_t part) {
    size_t res = length * part / 100;
    length -= res;
    return res;
}

std::string get_varname(std::mt19937 &gen) {
    std::uniform_int_distribution<int> len_dist(1, 16);
    std::uniform_int_distribution<int> letter_dist(0, 51);
    std::uniform_int_distribution<int> alnum_dist(0, 61);
    const std::string letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
    const std::string alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

    int len = len_dist(gen);
    std::string result;

    result += letters[letter_dist(gen)];

    for (int i = 1; i < len; ++i) {
        result += alnum[alnum_dist(gen)];
    }

    return result;
}

std::string get_varname_or_lit(std::mt19937 &gen, size_t &length) {
    if (gen() % 2 == 0) {
        return get_varname(gen);
    }
    size_t len = get_part(length, gen() % 50);
    constexpr char digits[] = "0123456789";
    std::string literal;
    literal.reserve(len);
    for (size_t i = 0; i < len; i++) {
        literal += digits[gen() % 10];
    }
    return literal;
}

std::string get_sign(size_t rnd) {
    constexpr char signs[] = "*/%";
    return {signs[rnd % 3]};
}

std::string gen_string(size_t length) {
    std::random_device random_device;
    std::mt19937 gen(random_device());
    std::stringstream ss;
    ss << get_vartype(gen()) << get_spaces(get_part(length, std::max(gen() % 30, static_cast<size_t>(1))))
            << get_varname(gen) << get_spaces(get_part(length, gen() % 20)) << "="
            << get_spaces(get_part(length, gen() % 20)) << get_varname_or_lit(gen, length);
    if (gen() % 2 == 0) {
        ss << get_spaces(length) << ";";
    } else {
        ss << get_spaces(get_part(length, gen() % 20)) << get_sign(gen()) << get_spaces(get_part(length, gen() % 20)) <<
                get_varname_or_lit(gen, length)
                << get_spaces(length) << ";";
    }


    return std::move(ss.str());
}

std::vector<std::string> generate_strings(size_t amount, size_t length) {
    std::vector<std::string> result;
    result.reserve(length);
    for (size_t i = 0; i < amount; i++) {
        result.emplace_back(gen_string(length));
    }
    return result;
}

int main() {
    Loader loader;
    std::vector<std::string> data = generate_strings(1, 5000);
    for (auto &i: loader.get_recognizers()) {
        auto start_time = std::chrono::high_resolution_clock::now();
        for (auto &s: data) {
            i.second.Validate(s);
        }
        auto end_time = std::chrono::high_resolution_clock::now();
        auto duration_us = std::chrono::duration_cast<std::chrono::microseconds>(end_time - start_time);

        std::cout << i.first << " execution time: " << duration_us.count() << " µs" << '\n';
    }
    return 0;
}

