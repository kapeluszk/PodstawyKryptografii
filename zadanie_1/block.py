import time
import helper
from Crypto.Cipher import AES

def mode_compare(key, data):
    choice = input("do you want to encrypt or decrypt? (e/d): ")
    if choice != "e" and choice != "d":
        print("Invalid input stopping the program")
        print(choice)
        return
    save = input("Do you want to save the encrypted files? (y/n): ")
    if save != "y" and save != "n":
        print("Invalid input stopping the program")
        print(save)
        return

    if choice == "e":
        start = time.time()
        encrypted = helper.encrypt_file(key, data, AES.MODE_ECB)
        end = time.time()
        print("Encryption time for ECB: ", end - start)
        if save == 'y':
            helper.save_file(encrypted, "ECB", choice)

        start = time.time()
        encrypted = helper.encrypt_file(key, data, AES.MODE_CBC)
        end = time.time()
        print("Encryption time for CBC: ", end - start)
        if save == 'y':
            helper.save_file(encrypted, "CBC", choice)

        start = time.time()
        encrypted = helper.encrypt_file(key, data, AES.MODE_CTR)
        end = time.time()
        print("Encryption time for CTR: ", end - start)
        if save == 'y':
            helper.save_file(encrypted, "CTR", choice)
    else:
        print("in which mode was the file encrypted? 1 - ecb/ 2 - cbc/ 3 - ctr")
        encryption_mode = int(input())
        if encryption_mode != 1 and encryption_mode != 2 and encryption_mode != 3:
            print("Invalid input stopping the program")
            print(encryption_mode)
            return
        if encryption_mode == 1:
            mode = AES.MODE_ECB
        elif encryption_mode == 2:
            mode = AES.MODE_CBC
        else:
            mode = AES.MODE_CTR

        start = time.time()
        decrypted = helper.decrypt_file(key, data, mode)
        end = time.time()
        print("Decryption time: ", end - start)
        if save == 'y':
            helper.save_file(decrypted, "ECB", choice)

def main():
    data = helper.select_file()
    key = b'0123456789abcdef'

    mode_compare(key, data)

if __name__ == "__main__":
    main()
