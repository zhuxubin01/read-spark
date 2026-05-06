import Combine
import Foundation

final class AuthRepository {
    static let shared = AuthRepository()
    private let api = APIService.shared

    func login(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        api.login(phone: phone, code: code)
            .handleEvents(receiveOutput: { tokens in
                KeychainManager.shared.accessToken = tokens.accessToken
                KeychainManager.shared.refreshToken = tokens.refreshToken
            })
            .eraseToAnyPublisher()
    }

    func register(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        api.register(phone: phone, code: code)
            .handleEvents(receiveOutput: { tokens in
                KeychainManager.shared.accessToken = tokens.accessToken
                KeychainManager.shared.refreshToken = tokens.refreshToken
            })
            .eraseToAnyPublisher()
    }

    func logout() {
        KeychainManager.shared.clearTokens()
    }
}
