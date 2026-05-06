import Combine
import Foundation

class ReaderViewModel: ObservableObject {
    @Published var article: Article?
    @Published var isLoading = false
    @Published var error: String?
    @Published var showTranslation = false

    private var cancellables = Set<AnyCancellable>()
    private let articleRepo = ArticleRepository.shared
    private let progressRepo = ProgressRepository.shared

    func loadArticle(id: UUID) {
        isLoading = true
        error = nil

        articleRepo.fetchArticle(id: id)
            .sink(receiveCompletion: { [weak self] completion in
                self?.isLoading = false
                if case .failure(let err) = completion {
                    self?.error = err.localizedDescription
                }
            }, receiveValue: { [weak self] article in
                self?.article = article
            })
            .store(in: &cancellables)
    }

    func syncProgress(position: Int, percentage: Double) {
        guard let article = article else { return }
        progressRepo.syncProgress(articleId: article.id, position: position, percentage: percentage)
            .sink(receiveCompletion: { _ in }, receiveValue: { _ in })
            .store(in: &cancellables)
    }
}
