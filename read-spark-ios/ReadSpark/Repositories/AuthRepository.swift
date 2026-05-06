import Combine
import Foundation

protocol AuthRepositoryProtocol {
    func login(phone: String, code: String) -> AnyPublisher<TokenPair, APIError>
    func register(phone: String, code: String) -> AnyPublisher<TokenPair, APIError>
    func logout()
}

final class AuthRepository {
    static let shared = AuthRepository()

    private let api: APIServiceProtocol
    private let tokenStore: TokenStore

    init(api: APIServiceProtocol = APIService.shared, tokenStore: TokenStore = KeychainManager.shared) {
        self.api = api
        self.tokenStore = tokenStore
    }

    func login(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        api.login(phone: phone, code: code)
            .handleEvents(receiveOutput: { tokens in
                self.tokenStore.accessToken = tokens.accessToken
                self.tokenStore.refreshToken = tokens.refreshToken
            })
            .eraseToAnyPublisher()
    }

    func register(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        api.register(phone: phone, code: code)
            .handleEvents(receiveOutput: { tokens in
                self.tokenStore.accessToken = tokens.accessToken
                self.tokenStore.refreshToken = tokens.refreshToken
            })
            .eraseToAnyPublisher()
    }

    func logout() {
        tokenStore.clearTokens()
    }
}

extension AuthRepository: AuthRepositoryProtocol {}
