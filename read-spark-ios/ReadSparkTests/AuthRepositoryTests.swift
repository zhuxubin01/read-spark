import Combine
import XCTest
@testable import ReadSpark

final class AuthRepositoryTests: XCTestCase {
    private var cancellables: Set<AnyCancellable> = []

    override func tearDown() {
        cancellables.removeAll()
        super.tearDown()
    }

    func testLoginStoresTokens() {
        let expected = TokenPair(accessToken: "acc-1", refreshToken: "ref-1", expiresIn: 900)
        let api = MockAPIService(loginResult: .success(expected))
        let tokenStore = MockTokenStore()
        let repo = AuthRepository(api: api, tokenStore: tokenStore)
        let exp = expectation(description: "login succeeds")

        repo.login(phone: "13800000000", code: "123456")
            .sink(receiveCompletion: { completion in
                if case .failure(let error) = completion {
                    XCTFail("unexpected failure: \(error)")
                }
            }, receiveValue: { value in
                XCTAssertEqual(value.accessToken, expected.accessToken)
                XCTAssertEqual(value.refreshToken, expected.refreshToken)
                XCTAssertEqual(tokenStore.accessToken, expected.accessToken)
                XCTAssertEqual(tokenStore.refreshToken, expected.refreshToken)
                exp.fulfill()
            })
            .store(in: &cancellables)

        wait(for: [exp], timeout: 1)
    }

    func testLogoutClearsTokens() {
        let tokenStore = MockTokenStore(accessToken: "old-a", refreshToken: "old-r")
        let repo = AuthRepository(api: MockAPIService(), tokenStore: tokenStore)
        repo.logout()

        XCTAssertNil(tokenStore.accessToken)
        XCTAssertNil(tokenStore.refreshToken)
        XCTAssertEqual(tokenStore.clearCount, 1)
    }
}

private final class MockTokenStore: TokenStore {
    var accessToken: String?
    var refreshToken: String?
    var clearCount = 0

    init(accessToken: String? = nil, refreshToken: String? = nil) {
        self.accessToken = accessToken
        self.refreshToken = refreshToken
    }

    func clearTokens() {
        accessToken = nil
        refreshToken = nil
        clearCount += 1
    }
}

private final class MockAPIService: APIServiceProtocol {
    var loginResult: Result<TokenPair, APIError>?
    var registerResult: Result<TokenPair, APIError>?
    var syncResult: Result<[String: String], APIError>?
    var syncArgs: (articleId: UUID, position: Int, percentage: Double)?

    init(loginResult: Result<TokenPair, APIError>? = nil,
         registerResult: Result<TokenPair, APIError>? = nil,
         syncResult: Result<[String: String], APIError>? = nil) {
        self.loginResult = loginResult
        self.registerResult = registerResult
        self.syncResult = syncResult
    }

    func login(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        switch loginResult ?? .failure(.serverError("missing login result")) {
        case .success(let tokens):
            return Just(tokens).setFailureType(to: APIError.self).eraseToAnyPublisher()
        case .failure(let error):
            return Fail(error: error).eraseToAnyPublisher()
        }
    }

    func register(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        switch registerResult ?? .failure(.serverError("missing register result")) {
        case .success(let tokens):
            return Just(tokens).setFailureType(to: APIError.self).eraseToAnyPublisher()
        case .failure(let error):
            return Fail(error: error).eraseToAnyPublisher()
        }
    }

    func syncProgress(articleId: UUID, position: Int, percentage: Double) -> AnyPublisher<[String: String], APIError> {
        syncArgs = (articleId, position, percentage)
        switch syncResult ?? .success(["status": "ok"]) {
        case .success(let payload):
            return Just(payload).setFailureType(to: APIError.self).eraseToAnyPublisher()
        case .failure(let error):
            return Fail(error: error).eraseToAnyPublisher()
        }
    }
}
