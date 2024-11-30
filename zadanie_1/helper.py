import os

from Crypto.Cipher import AES
from Crypto.Util.Padding import pad, unpad

def select_file():
    print("Select file:")
    filepath = input()
    with open(filepath, 'rb') as file:
        data = file.read()
    return data

def save_file(data, mode, choice):
    if choice == "e":
        filepath = "encrypted_" + mode + ".txt"
    else:
        filepath = "decrypted_" + mode + ".txt"
    with open(filepath, 'wb') as file:
        file.write(data)

def encrypt_file(key, data, mode):
    cipher = AES.new(key, mode)
    return cipher.encrypt(data)

def decrypt_file(key, data, mode):
    cipher = AES.new(key, mode, iv=data[:16] if mode == AES.MODE_CBC else None)
    decrypted_data = cipher.decrypt(data[16:] if mode == AES.MODE_CBC else data)
    return unpad(decrypted_data, 16)

def my_own_cbc_encrypt(key, data):
    iv = os.urandom(16)  # generate random initialization vector
    encrypted = iv
    previous_block = iv

    data = pad(data, AES.block_size)  # ensure that the data is multiple of 16 bytes

    for i in range(0, len(data), 16):
        block = data[i:i+16]
        block = bytes(a ^ b for a, b in zip(block, previous_block))  # XOR with the previous block
        cipher = AES.new(key, AES.MODE_ECB)
        encrypted_block = cipher.encrypt(block)
        encrypted += encrypted_block
        previous_block = encrypted_block

    return encrypted

def my_own_cbc_decrypt(key, data):
    iv = data[:16]
    decrypted = b''
    previous_block = iv

    for i in range(16, len(data), 16):
        block = data[i:i+16]
        cipher = AES.new(key, AES.MODE_ECB)
        decrypted_block = cipher.decrypt(block)
        decrypted_block = bytes(a ^ b for a, b in zip(decrypted_block, previous_block))  # XOR with the previous block
        decrypted += decrypted_block
        previous_block = block

    return decrypted