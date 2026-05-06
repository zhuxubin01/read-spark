import Combine
import Foundation

protocol ProgressRepositoryProtocol {
    func syncProgress(articleId: UUID, position: Int, percentage: Double) -> AnyPublisher<Void, APIError>
}

final class ProgressRepository {
    static let shared = ProgressRepository()

    private let api: APIServiceProtocol

    init(api: APIServiceProtocol = APIService.shared) {
        self.api = api
    }

    func syncProgress(articleId: UUID, position: Int, percentage: Double) -> AnyPublisher<Void, APIError> {
        api.syncProgress(articleId: articleId, position: position, percentage: percentage)
            .map { _ in () }
            .eraseToAnyPublisher()
    }
}

extension ProgressRepository: ProgressRepositoryProtocol {}
