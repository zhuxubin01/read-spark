import Combine
import XCTest
@testable import ReadSpark

final class ProgressRepositoryTests: XCTestCase {
    private var cancellables: Set<AnyCancellable> = []

    override func tearDown() {
        cancellables.removeAll()
        super.tearDown()
    }

    func testSyncProgressMapsResponseToVoidAndForwardsParameters() {
        let articleId = UUID()
        let api = MockAPIService(syncResult: .success(["status": "ok"]))
        let repo = ProgressRepository(api: api)
        let exp = expectation(description: "sync succeeds")

        repo.syncProgress(articleId: articleId, position: 120, percentage: 0.42)
            .sink(receiveCompletion: { completion in
                if case .failure(let error) = completion {
                    XCTFail("unexpected failure: \(error)")
                }
            }, receiveValue: { _ in
                exp.fulfill()
            })
            .store(in: &cancellables)

        wait(for: [exp], timeout: 1)
        XCTAssertEqual(api.syncArgs?.articleId, articleId)
        XCTAssertEqual(api.syncArgs?.position, 120)
        guard let percentage = api.syncArgs?.percentage else {
            return XCTFail("missing sync percentage argument")
        }
        XCTAssertEqual(percentage, 0.42, accuracy: 0.0001)
    }
}

private final class MockAPIService: APIServiceProtocol {
    var syncResult: Result<[String: String], APIError>?
    var syncArgs: (articleId: UUID, position: Int, percentage: Double)?

    init(syncResult: Result<[String: String], APIError>? = nil) {
        self.syncResult = syncResult
    }

    func login(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        Fail(error: .serverError("unused")).eraseToAnyPublisher()
    }

    func register(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        Fail(error: .serverError("unused")).eraseToAnyPublisher()
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
