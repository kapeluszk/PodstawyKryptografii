import math
import random
import sympy


def my_own_rsa_encrypt(message, e, n):
    return pow(message, e, n)

def my_own_rsa_decrypt(cipher, d, n):
    return pow(cipher, d, n)

def generate_x(N):
    while True:
        x = random.randint(1,N-1)
        if math.gcd(x,N) == 1:
            break
    return x

def generate_big_primes():
    same = True
    while same:
        p = sympy.nextprime(random.randint(100000000000000000, 1000000000000000000))
        q = sympy.nextprime(random.randint(100000000000000000, 1000000000000000000))
        if p != q:
            same = False
    return p, q

def my_own_rsa_keygen(p, q):
    n = p * q
    phi = (p - 1) * (q - 1)
    e = generate_x(phi)
    d = pow(e, -1, phi)
    return e, d, n

def main():
    p,q = generate_big_primes()
    e, d, n = my_own_rsa_keygen(p, q)
    message = 1234567890
    print("Message: ", message)
    cipher = my_own_rsa_encrypt(message, e, n)
    print("Cipher: ", cipher)
    decrypted = my_own_rsa_decrypt(cipher, d, n)
    print("Decrypted: ", decrypted)



if __name__ == "__main__":
    main()