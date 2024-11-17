def bit_test(bit_output, n):
    one_count = 0
    for i in range(n):
        if bit_output[i] == '1':
            one_count += 1

    min = 9725
    max = 10275
    if one_count < min or one_count > max:
        return False
    else:
        return True

def long_series_test(bit_output, n):
    for i in range(n - 26):
        if bit_output[i:i + 26] == '1' * 26:
            return False
        elif bit_output[i:i + 26] == '0' * 26:
            return False
    return True

def series_test(bit_output, n):
    ranges_dict = {
        1: (2315, 2685),
        2: (1114, 1386),
        3: (527, 723),
        4: (240, 384),
        5: (103, 209),
        6: (103, 209),
    }

    def count_series(bit_output, n, bit_value):
        amount_of_series = {i: 0 for i in range(1, 7)}
        i = 0
        while i < n:
            count = 1
            while i + count < n and bit_output[i] == bit_value and bit_output[i] == bit_output[i + count]:
                count += 1

            if bit_output[i] == bit_value:
                if count >= 6:
                    amount_of_series[6] += 1
                else:
                    amount_of_series[count] += 1

            i += count

        for length, (min_count, max_count) in ranges_dict.items():
            if amount_of_series[length] < min_count or amount_of_series[length] > max_count:
                return False
        return True

    result_ones = count_series(bit_output, n, '1')
    result_zeros = count_series(bit_output, n, '0')

    if result_ones and result_zeros:
        return True
    else:
        return False

def poker_test(bit_output, n):
    num_groups = 5000
    counts = {format(i, '04b'): 0 for i in range(16)}

    for i in range(num_groups):
        group = bit_output[i * 4:(i + 1) * 4]
        counts[group] += 1

    x = (16 / num_groups) * sum(count ** 2 for count in counts.values()) - num_groups

    min_x = 2.16
    max_x = 46.17

    return min_x <= x <= max_x

