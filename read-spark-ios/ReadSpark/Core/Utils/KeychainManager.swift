import Foundation

final class KeychainManager {
    static let shared = KeychainManager()

    private var storage: [String: String] = [:]

    var accessToken: String? {
        storage["access_token"]
    }

    private init() {}
}
