import math
import random
import sympy
import tests

def generate_prime(value):
    p = sympy.nextprime(value)
    while p % 4 != 3:
        p = sympy.nextprime(p)
    return p

def generate_x(N):
    while True:
        x = random.randint(1,N-1)
        if math.gcd(x,N) == 1:
            break
    return x

def generate_bits(n):
    big_num1 = 3 * random.randint(1000000, 10000000)
    big_num2 = 4 * random.randint(1000000, 10000000)
    p = generate_prime(big_num1)
    q = generate_prime(big_num2)
    N = p * q
    x = generate_x(N)

    x0 = x * x % N

    bit_output = ''
    xi = x0
    for _ in range(n):
        xi = xi * xi % N
        b = xi % 2
        bit_output += str(b)
    return bit_output

def main():
    n = 20000
    while True:
        bit_output = generate_bits(n)
        if tests.poker_test(bit_output, n) and tests.series_test(bit_output, n) and tests.bit_test(bit_output, n) and tests.long_series_test(bit_output, n):
            break

    # print("p:",p)
    # print("q:",q)
    # print("N:",N)
    # print("x:",x)
    # print("x0:",x0,"\n")
    # print("Did pass bit test: ", tests.bit_test(bit_output, n))
    # print("Did pass long series test: ", tests.long_series_test(bit_output, n))
    # print("Did pass series test: ", tests.series_test(bit_output, n))
    # print("Did pass poker test: ", tests.poker_test(bit_output, n),"\n")


    print("generated bits passed all tests! \nbits: ",bit_output)


if __name__ == "__main__":
    main()



