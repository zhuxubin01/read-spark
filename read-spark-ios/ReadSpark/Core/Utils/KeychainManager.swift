import Foundation
import Security

protocol TokenStore: AnyObject {
    var accessToken: String? { get set }
    var refreshToken: String? { get set }
    func clearTokens()
}

final class KeychainManager {
    static let shared = KeychainManager()

    private let service = "com.readspark.app"

    var accessToken: String? {
        get { get(key: "accessToken") }
        set { set(key: "accessToken", value: newValue) }
    }

    var refreshToken: String? {
        get { get(key: "refreshToken") }
        set { set(key: "refreshToken", value: newValue) }
    }

    private init() {}

    private func get(key: String) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecAttrService as String: service,
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne
        ]

        var result: AnyObject?
        SecItemCopyMatching(query as CFDictionary, &result)

        guard let data = result as? Data else { return nil }
        return String(data: data, encoding: .utf8)
    }

    private func set(key: String, value: String?) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecAttrService as String: service
        ]

        SecItemDelete(query as CFDictionary)

        guard let value, let data = value.data(using: .utf8) else { return }

        let attributes: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecAttrService as String: service,
            kSecValueData as String: data
        ]

        SecItemAdd(attributes as CFDictionary, nil)
    }

    func clearTokens() {
        set(key: "accessToken", value: nil)
        set(key: "refreshToken", value: nil)
    }
}

extension KeychainManager: TokenStore {}
